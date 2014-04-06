package angularjs

import (
	"reflect"

	"github.com/codegangsta/inject"
	"github.com/gopherjs/gopherjs/js"
)

var (
	Ng = InitAngular()
)

type DummyJsObj struct{ js.Object }

func (obj DummyJsObj) Get(name string) js.Object {
	return obj
}

func (obj DummyJsObj) Call(name string, params ...interface{}) js.Object {
	return obj
}

func AngularJs() (obj js.Object) {
	defer func() {
		if r := recover(); r != nil {
			obj = DummyJsObj{}
		}
	}()
	return js.Global.Get("angular")
}

type Injector struct {
	inject.Injector
	providers map[reflect.Type]provider
}

func NewInjector() *Injector {
	inj := &Injector{
		inject.New(),
		map[reflect.Type]provider{
			reflect.TypeOf(RouteProvider{}): &RouteProvider{NewProvider("$routeProvider")},
		},
	}

	return inj
}

//RequestedProviders gets the list of dependencies required in the funtion fn's
//parameter list.
func (inj *Injector) RequestedProviders(fn interface{}) (providers []provider) {
	t := reflect.TypeOf(fn)
	providers = make([]provider, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		argType := t.In(i).Elem()
		var ok bool
		providers[i], ok = inj.providers[argType]
		if !ok {
			panic("Invalid provider.")
		}
	}
	return
}

//AngularDeps makes a generic array of providers
//and add the injected function fn at the end.
func (inj *Injector) AngularDeps(fn interface{}) []interface{} {
	rp := inj.RequestedProviders(fn)
	deps := make([]interface{}, len(rp))
	for i, _ := range rp {
		deps[i] = rp[i].AngularName()
	}
	deps = append(deps, func(providers ...js.Object) {
		in := make([]reflect.Value, len(rp))
		for i, p := range rp {
			p.SetJs(providers[i])
			in[i] = reflect.ValueOf(p)
		}
		reflect.ValueOf(fn).Call(in)
	})
	return deps
}

type Angular struct {
	Inj *Injector //dependency injector
}

func InitAngular() *Angular {
	ng := &Angular{NewInjector()}
	return ng
}

type Module struct{ js.Object }
type RouteProvider struct {
	*Provider
}

type provider interface {
	SetJs(js.Object)
	AngularName() string
}

type Provider struct {
	js.Object
	angularName string
}

func NewProvider(name string) *Provider {
	p := &Provider{}
	p.angularName = name
	return p
}

func (p *Provider) SetJs(obj js.Object) {
	p.Object = obj
}

func (p *Provider) AngularName() string {
	return p.angularName
}

func (r *RouteProvider) When(path string, route *AngularOpts) *RouteProvider {
	r.Call("when", path, route.toJs())
	return r
}

func (r *RouteProvider) Otherwise(route *AngularOpts) {
	r.Call("otherwise", route.toJs())
}

func (m *Module) Config(fn interface{}) {
	m.Call("config", Ng.Inj.AngularDeps(fn))
}

func (m *Module) NewController(name string, constructor func(scope *Scope)) {
	m.Call("controller", name, func(dollar_scope js.Object) {
		constructor(&Scope{dollar_scope})
	})
}

type Scope struct{ js.Object }

func (s *Scope) Apply(f func()) {
	s.Call("$apply", f)
}

func (s *Scope) EvalAsync(f func()) {
	s.Call("$evalAsync", f)
}

type JQueryElement struct{ js.Object }

func (e *JQueryElement) Prop(name string) js.Object {
	return e.Call("prop", name)
}

func (e *JQueryElement) SetProp(name, value interface{}) {
	e.Call("prop", name, value)
}

func (e *JQueryElement) On(events string, handler func(*Event)) {
	e.Call("on", events, func(e js.Object) {
		handler(&Event{Object: e})
	})
}

func (e *JQueryElement) Val() js.Object {
	return e.Call("val")
}

func (e *JQueryElement) SetVal(value interface{}) {
	e.Call("val", value)
}

type Event struct {
	js.Object
	KeyCode int `js:"keyCode"`
}

func (e *Event) PreventDefault() {
	e.Call("preventDefault")
}

//Bootstrap triggers angular's bootstrap
func Bootstrap(modules []string) {
	AngularJs().Call("bootstrap",
		js.Global.Get("document"), modules)
}

func NewModule(name string, requires []string) *Module {
	return &Module{AngularJs().Call("module", name, requires)}
}

func ElementById(id string) *JQueryElement {
	return &JQueryElement{AngularJs().Call("element", js.Global.Get("document").Call("getElementById", id))}
}

func Service(name string) js.Object {
	return AngularJs().Call("element", js.Global.Get("document")).Call("injector").Call("get", name)
}

type HttpService struct{}

type HttpMethod string

var HTTP = new(HttpService)

const (
	HttpGet  HttpMethod = "GET"
	HttpPost HttpMethod = "POST"
)

type Future struct{ js.Object }

type RequestCallback func(data string, status int)

func (ft *Future) call(method string, callback RequestCallback) *Future {
	ft.Object.Call(method, func(data string, status int, headers js.Object, config js.Object) {
		callback(data, status)
	})
	return ft
}

func (ft *Future) Success(callback RequestCallback) *Future {
	return ft.call("success", callback)
}
func (ft *Future) Error(callback RequestCallback) *Future {
	return ft.call("error", callback)
}

//Req performs a http request
func (s *HttpService) Req(method HttpMethod, url string) *Future {
	future := Service("$http").Invoke(map[string]string{
		"method": string(method),
		"url":    url,
	})
	return &Future{future}
}

func (s *HttpService) Get(url string) *Future {
	return s.Req(HttpGet, url)
}
