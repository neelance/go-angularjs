package angularjs

import (
	"reflect"

	"github.com/gopherjs/gopherjs/js"
	"github.com/pkg/errors"
)

var (
	injectMap = make(map[reflect.Type]*Inject)
)

// Inject represents a resource that can be injected into an controller/directive/route-resolve func
type Inject struct {
	f           func(obj *js.Object) reflect.Value
	angularName string
}

// RegisterResource registers a resource that can be injected
func RegisterResource(resourceType reflect.Type, angularName string, f func(obj *js.Object) reflect.Value) {
	injectMap[resourceType] = &Inject{
		f:           f,
		angularName: angularName,
	}
}

// GetResource returns the resource for a given reflect.Type
func GetResource(resourceType reflect.Type) *Inject {
	return injectMap[resourceType]
}

// MakeFuncInjectable returns a func that transforms *js.Object's to
// the corresponding go types.
func MakeFuncInjectable(f interface{}) (jsFunc js.S, err error) {
	angularParamNames := make(js.S, 0)
	transFormFuncs := make([]func(obj *js.Object) reflect.Value, 0)

	injects, callable, err := GetFuncInjectables(f)
	if err != nil {
		return nil, err
	}

	for _, i := range injects {
		angularParamNames = append(angularParamNames, i.angularName)
		transFormFuncs = append(transFormFuncs, i.f)
	}

	// we return an optional *js.Object here if the function returns it
	return append(angularParamNames, func(objs ...*js.Object) *js.Object {
		args := make([]reflect.Value, 0)
		for i, obj := range objs {
			args = append(args, transFormFuncs[i](obj))
		}

		ret := callable.Call(args)
		if len(ret) == 0 {
			return nil
		}

		return ret[0].Interface().(*js.Object)
	}), nil
}

// GetFuncInjectables returns the reflect.Type's resources that have to be uses
// to call the given function correctly.
func GetFuncInjectables(f interface{}) (injects []*Inject, callable reflect.Value, err error) {
	callable = reflect.ValueOf(f)
	if callable.Kind() != reflect.Func {
		return nil, callable, errors.New("Only func's can be made injectable")
	}

	for i := 0; i < callable.Type().NumIn(); i++ {
		arg := callable.Type().In(i)

		if injector := GetResource(arg); injector != nil {
			injects = append(injects, injector)
		} else {
			return nil, callable, errors.Errorf("no resource found for type: %v", arg)
		}
	}

	return injects, callable, nil
}
