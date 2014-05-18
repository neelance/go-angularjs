package angularjs

import (
	"fmt"
	"reflect"

	"github.com/codegangsta/inject"
	"github.com/gopherjs/gopherjs/js"
	op "github.com/phaikawl/options"
)

var (
	Ng                   = InitAngular()
	AngularProvidersList = map[reflect.Type]provider{
		reflect.TypeOf(RouteProvider{}): &RouteProvider{NewProvider("$routeProvider")},
		reflect.TypeOf(HttpService{}):   &HttpService{NewProvider("$http")},
		reflect.TypeOf(Scope{}):         &Scope{NewProvider("$scope")},
		reflect.TypeOf(RootScope{}):     &RootScope{NewProvider("$rootScope")},
		reflect.TypeOf(QProvider{}):     &QProvider{NewProvider("$q")},
		reflect.TypeOf(HttpProvider{}):  &HttpProvider{NewProvider("$httpProvider")},
		reflect.TypeOf(Interval{}):      &Interval{NewProvider("$interval")},
	}
)

type AngularOpts struct{ *op.OptionsProvider }

func (o *AngularOpts) toJs() map[string]interface{} {
	return o.OptionsProvider.ExportToMapWithTag("js")
}

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
		AngularProvidersList,
	}

	return inj
}

//requestedProviders gets the list of dependencies required in the funtion fn's
//parameter list.
func (inj *Injector) requestedProviders(fn interface{}) (providers []provider) {
	t := reflect.TypeOf(fn)
	providers = make([]provider, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		argType := t.In(i).Elem()
		var ok bool
		providers[i], ok = inj.providers[argType]
		if !ok {
			panic(fmt.Sprintf("Invalid provider type %v.", argType.Name()))
		}
	}
	return
}

func (inj *Injector) _angularDeps(fn interface{}, transformFn func(reflect.Value) reflect.Value) []interface{} {
	rp := inj.requestedProviders(fn)
	deps := make([]interface{}, len(rp))
	for i, _ := range rp {
		deps[i] = rp[i].AngularName()
	}
	deps = append(deps, func(providers ...js.Object) interface{} {
		in := make([]reflect.Value, len(rp))
		for i, p := range rp {
			p.SetJs(providers[i])
			in[i] = reflect.ValueOf(p)
		}
		if reflect.TypeOf(fn).NumOut() > 0 && transformFn != nil {
			return transformFn(reflect.ValueOf(fn).Call(in)[0]).Interface()
		}
		reflect.ValueOf(fn).Call(in)
		return nil
	})
	return deps
}

//angularDeps makes a generic array of providers
//and add the injected function fn at the end.
func (inj *Injector) angularDeps(fn interface{}) []interface{} {
	return inj._angularDeps(fn, func(v reflect.Value) reflect.Value {
		return v
	})
}

type Angular struct {
	Inj *Injector //dependency injector
}

func InitAngular() *Angular {
	ng := &Angular{NewInjector()}
	return ng
}

type Module struct{ js.Object }

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

func (m *Module) Config(fn interface{}) {
	m.Call("config", Ng.Inj.angularDeps(fn))
}

func (m *Module) Factory(name string, fn interface{}) {
	m.Call("factory", name, Ng.Inj.angularDeps(fn))
}

func (m *Module) NewController(name string, constructor interface{}) {
	m.Call("controller", name, Ng.Inj.angularDeps(constructor))
}

type Interval struct{ *Provider }

type Scope struct{ *Provider }

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
