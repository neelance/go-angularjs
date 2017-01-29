package angularjs

import (
	"reflect"

	"github.com/gopherjs/gopherjs/js"
)

func init() {

	RegisterResource(reflect.TypeOf(&Scope{}), "$scope", func(obj *js.Object) reflect.Value {
		return reflect.ValueOf(&Scope{Object: obj})
	})

	RegisterResource(reflect.TypeOf(&RootScope{}), "$rootScope", func(obj *js.Object) reflect.Value {
		return reflect.ValueOf(&RootScope{Object: obj})
	})

	RegisterResource(reflect.TypeOf(&Q{}), "$q", func(obj *js.Object) reflect.Value {
		return reflect.ValueOf(&Q{Object: obj})
	})

}

// NewModule makes a new angular module
func NewModule(name string, requires []string, configFn interface{}) *Module {
	transformedFunc, err := MakeFuncInjectable(configFn)
	if err != nil {
		panic(err)
	}
	return &Module{js.Global.Get("angular").Call("module", name, requires, transformedFunc)}
}

// Service gets a service trough the global $injector
func Service(name string) *js.Object {
	return js.Global.Get("angular").Call("element", js.Global.Get("document")).Call("injector").Call("get", name)
}

// Module is an angular module
type Module struct{ *js.Object }

// NewController creates a new controller on the given module
// The arguments of the controller should only contain resources that have
// been added with a call to RegisterResouce
func (m *Module) NewController(name string, constructor interface{}) {
	transformedFunc, err := MakeFuncInjectable(constructor)
	if err != nil {
		panic(err)
	}
	m.Call("controller", name, transformedFunc)
}

// NewDirective registers a directive.
// The fomat should be: func(.... resources) Directive
func (m *Module) NewDirective(name string, f interface{}) {
	injects, callable, err := GetFuncInjectables(f)
	if err != nil {
		panic(err)
	}

	if callable.Type().NumOut() != 1 {
		panic("expected 1 out arg for directive: " + name)
	}

	if callable.Type().Out(0) != reflect.TypeOf(Directive{}) {
		panic("directive: " + name + " should have 1 return value of type: " + reflect.TypeOf(Directive{}).String() + " got: " + callable.Type().Out(0).String())
	}

	angularParamNames := make(js.S, 0)
	transFormFuncs := make([]func(obj *js.Object) reflect.Value, 0)
	for _, i := range injects {
		angularParamNames = append(angularParamNames, i.angularName)
		transFormFuncs = append(transFormFuncs, i.f)
	}

	m.Call("directive", []interface{}{name, append(angularParamNames, func(objs ...*js.Object) map[string]interface{} {
		args := make([]reflect.Value, 0)
		for i, obj := range objs {
			args = append(args, transFormFuncs[i](obj))
		}

		ret := callable.Call(args)
		if len(ret) == 0 {
			return nil
		}

		directive := ret[0].Interface().(Directive)
		directiveMap := make(map[string]interface{})
		directiveMap["templateUrl"] = directive.TemplateURL
		scope := make(map[string]string)
		for _, bindAttr := range directive.Scope {
			scope[bindAttr] = "="
		}

		directiveMap["scope"] = scope
		directiveMap["restrict"] = directive.Restrict

		if directive.Link != nil {
			directiveMap["link"] = func(scope *js.Object, element *js.Object, attrs map[string]interface{}, controller *js.Object, transcludeFn *js.Object) {
				directive.Link(&Scope{scope}, &Element{element}, attrs)
			}
		}

		return directiveMap
	})}...)
}

// Directive is the go variant of an angular directive
// At least a templateURL should be given.
type Directive struct {
	Restrict    string
	Link        func(scope *Scope, el *Element, attrs map[string]interface{})
	TemplateURL string
	Scope       []string
}

// RootScope is the angular $rootScope
type RootScope Scope

// Scope is the angular $scope
type Scope struct {
	*js.Object
}

// Apply is the angular $scope.$apply and should be called when
// changing values trough goroutines or event listeners
func (s *Scope) Apply(f func()) {
	s.Call("$apply", f)
}

// CopyScope is used to provide a *js.Object scope for structs inside a controller
func (s *Scope) CopyScope() *js.Object {
	return s.Call("$new")
}

// Watch is the angular $watch
func (s *Scope) Watch(key string, f func(newValue interface{}, oldValue interface{})) {
	s.Call("$watch", key, f)
}

func (s *Scope) EvalAsync(f func()) {
	s.Call("$evalAsync", f)
}

// On registers the given listener for an event.
// Returns an unregister func for the event
// format of listener: func(ev *Event, args ...interface{})
func (s *Scope) On(event string, listener interface{}) (unregisterFunc func()) {
	l := s.Call("$on", event, listener)

	return func() {
		l.Invoke()
	}
}

// Listen binds to the event and removes it if the scope is destroyed
func (s *Scope) Listen(event string, listener interface{}) {
	l := s.On(event, listener)
	s.On("$destroy", l)
}

// Emit is a copy of $emit
func (s *Scope) Emit(event string, args ...interface{}) {
	args = append([]interface{}{event}, args...)
	s.Call("$emit", args...)
}

// Broadcast is a copy of $broadcast
func (s *Scope) Broadcast(event string, args ...interface{}) {
	args = append([]interface{}{event}, args...)
	s.Call("$broadcast", args...)
}

// Event is an angular event from $emit/$broadcast
type Event struct {
	*js.Object
}

// TargetScope see $event.targetScope
func (e *Event) TargetScope() *Scope {
	return &Scope{e.Get("targetScope")}
}

// CurrentScope see $event.currentScope
func (e *Event) CurrentScope() *Scope {
	return &Scope{e.Get("currentScope")}
}

// Name see $event.name
func (e *Event) Name() string {
	return e.Get("name").String()
}

// PreventDefault see $event.preventDefault()
func (e *Event) PreventDefault() {
	e.Call("preventDefault")
}

// DefaultPrevented see $event.defaultPrevented
func (e *Event) DefaultPrevented() bool {
	return e.Get("defaultPrevented").Bool()
}

// Q represets angulars $q
type Q struct {
	*js.Object
}

func (q *Q) Defer() *Deferred {
	return &Deferred{q.Call("defer")}
}

// Deferred is angulars defer object returned from $q.defer()
type Deferred struct {
	*js.Object
}

// Promise is angulars deferred.promis
func (d *Deferred) Promise() *js.Object {
	return d.Get("promise")
}

// Resolve will resolve the promise.
// Only js objects with an initialized object can be given.
// An empty javascript object can be created with 'js.Global.Get("Object").New()'
func (d *Deferred) Resolve(data interface{}) *js.Object {
	return d.Call("resolve", data)
}

// Reject see Resolve
func (d *Deferred) Reject(data interface{}) *js.Object {
	return d.Call("reject", data)
}
