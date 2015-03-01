package angularjs

import "github.com/gopherjs/gopherjs/js"

type Module struct{ *js.Object }

func (m *Module) NewController(name string, constructor func(scope *Scope)) {
	m.Call("controller", name, js.S{"$scope", func(scope *js.Object) {
		constructor(&Scope{scope})
	}})
}

type Scope struct{ *js.Object }

func (s *Scope) Apply(f func()) {
	s.Call("$apply", f)
}

func (s *Scope) EvalAsync(f func()) {
	s.Call("$evalAsync", f)
}

type JQueryElement struct{ *js.Object }

func (e *JQueryElement) Prop(name string) *js.Object {
	return e.Call("prop", name)
}

func (e *JQueryElement) SetProp(name, value interface{}) {
	e.Call("prop", name, value)
}

func (e *JQueryElement) On(events string, handler func(*Event)) {
	e.Call("on", events, func(e *js.Object) {
		handler(&Event{Object: e})
	})
}

func (e *JQueryElement) Val() *js.Object {
	return e.Call("val")
}

func (e *JQueryElement) SetVal(value interface{}) {
	e.Call("val", value)
}

type Event struct {
	*js.Object
	KeyCode int `js:"keyCode"`
}

func (e *Event) PreventDefault() {
	e.Call("preventDefault")
}

func NewModule(name string, requires []string, configFn func()) *Module {
	return &Module{js.Global.Get("angular").Call("module", name, requires, configFn)}
}

func ElementById(id string) *JQueryElement {
	return &JQueryElement{js.Global.Get("angular").Call("element", js.Global.Get("document").Call("getElementById", id))}
}

func Service(name string) *js.Object {
	return js.Global.Get("angular").Call("element", js.Global.Get("document")).Call("injector").Call("get", name)
}

type HttpService struct{}

var HTTP = new(HttpService)

func (s *HttpService) Get(url string, callback func(data string, status int)) {
	future := Service("$http").Call("get", url)
	future.Call("success", func(data string, status int, headers *js.Object, config *js.Object) {
		callback(data, status)
	})
	future.Call("error", func(data string, status int, headers *js.Object, config *js.Object) {
		callback(data, status)
	})
}
