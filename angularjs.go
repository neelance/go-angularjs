package angularjs

import (
	"fmt"

	"github.com/gopherjs/gopherjs/js"
)

type Module struct{ js.Object }
type routeProvider struct{ js.Object }

func (r *routeProvider) When(path string, route *ConfigMap) *routeProvider {
	r.Call("when", path, route.Options())
	return r
}

func (r *routeProvider) Otherwise(route *ConfigMap) {
	r.Call("otherwise", route.Options())
}

//depList makes a generic array of providers
//and add the function fn at the end.
func depList(fn interface{}, providers []string) []interface{} {
	l := make([]interface{}, len(providers)+1)
	for i, p := range providers {
		l[i] = p
	}
	l[len(providers)] = fn
	return l
}

type Injector struct {
	providerNames []string
	providers     map[string]js.Object
}

func Inject(providers ...string) *Injector {
	return &Injector{providers, make(map[string]js.Object)}
}

func (inj *Injector) GetProvider(name string) js.Object {
	if p, ok := inj.providers[name]; ok {
		return p
	}
	panic(fmt.Sprintf("Provider %s hasn't been injected.", name))
}

func (inj *Injector) injectReal(providers []js.Object) {
	for i, name := range inj.providerNames {
		inj.providers[name] = providers[i]
	}
}

type injectedFunc func(inj *Injector)

//RouteProvider retrieves the angular $routeProvider from the Injector
func RouteProvider(inj *Injector) *routeProvider {
	return &routeProvider{inj.GetProvider("$routeProvider")}
}

func (m *Module) Config(fn injectedFunc, inj *Injector) {
	m.Call("config", depList(func(providers ...js.Object) {
		inj.injectReal(providers)
		fn(inj)
	}, inj.providerNames))
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
	js.Global.Get("angular").Call("bootstrap",
		js.Global.Get("document"), modules)
}

func NewModule(name string, requires []string) *Module {
	return &Module{js.Global.Get("angular").Call("module", name, requires)}
}

func ElementById(id string) *JQueryElement {
	return &JQueryElement{js.Global.Get("angular").Call("element", js.Global.Get("document").Call("getElementById", id))}
}

func Service(name string) js.Object {
	return js.Global.Get("angular").Call("element", js.Global.Get("document")).Call("injector").Call("get", name)
}

type HttpService struct{}

var HTTP = new(HttpService)

func (s *HttpService) Get(url string, callback func(data string, status int)) {
	future := Service("$http").Call("get", url)
	future.Call("success", func(data string, status int, headers js.Object, config js.Object) {
		callback(data, status)
	})
	future.Call("error", func(data string, status int, headers js.Object, config js.Object) {
		callback(data, status)
	})
}
