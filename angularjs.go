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
	providers map[reflect.Type]string
}

func NewInjector() *Injector {
	inj := &Injector{
		inject.New(),
		map[reflect.Type]string{
			reflect.TypeOf(&RouteProvider{}): "$routeProvider",
		},
	}

	//for typ, _ := range DefaultProvidersMap {
	//	rtype := reflect.TypeOf(typ)
	//	v := reflect.ValueOf(typ)
	//	inj.Set(rtype, v)
	//}

	return inj
}

func (inj *Injector) RequestedProviders(fn interface{}) (types []reflect.Type, names []string) {
	t := reflect.TypeOf(fn)
	names = make([]string, t.NumIn())
	types = make([]reflect.Type, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		argType := t.In(i)
		names[i], types[i] = inj.providers[argType], argType
	}
	return
}

////depList makes a generic array of providers
////and add the function fn at the end.
//func (inj *Injector) DepList(fn interface{}) []interface{} {
//	providers := inj.InjectedProviders()
//	l := make([]interface{}, len(providers)+1)
//	for i, p := range providers {
//		l[i] = fmt.Sprintf("%sProvider", p)
//	}
//	l[len(providers)] = fn
//	return l
//}

func (inj *Injector) AngularDeps(fn interface{}) []interface{} {
	ptypes, pnames := inj.RequestedProviders(fn)
	deps := make([]interface{}, len(pnames))
	for i, _ := range pnames {
		deps[i] = pnames[i]
	}
	deps = append(deps, func(providers ...js.Object) {
		in := make([]reflect.Value, len(ptypes))
		for i, ptype := range ptypes {
			in[i] = reflect.New(ptype.Elem())
			if !in[i].Elem().FieldByName("Object").CanSet() {
				panic("Something's wrong with the provider object.")
			}
			field := in[i].Elem().FieldByName("Object")
			field.Set(reflect.ValueOf(providers[i]))
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
	js.Object
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
