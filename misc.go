package angularjs

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"reflect"
)

func returnTypeMustBe(typeName string, fn interface{}) {
	fnt := reflect.TypeOf(fn)
	if fnt.NumOut() != 1 || fnt.Out(0).Name() != typeName {
		panic(fmt.Sprintf("wrong return type %v for the function, expected %v", fnt.Out(0).Name(), typeName))
	}
}

type RootScope Scope

type QProvider struct{ *Provider }

func (q *QProvider) Defer() js.Object {
	return q.Object.Call("defer")
}

func (q *QProvider) When(value js.Object) js.Object {
	return q.Object.Call("when", value)
}

func (q *QProvider) Reject(reason string) js.Object {
	return q.Object.Call("reject", reason)
}

func (q *QProvider) All(promises []js.Object) js.Object {
	return q.Object.Call("all", promises)
}

func (q *QProvider) NowOrLater(obj js.Object) interface{} {
	if obj.IsNull() || obj.IsUndefined() {
		return q.When(obj)
	}
	return obj
}
