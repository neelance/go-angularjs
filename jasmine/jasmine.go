package jasmine

import "github.com/gopherjs/gopherjs/js"

type Expectation struct{ js.Object }

func (e *Expectation) ToBe(value interface{}) {
	e.Object.Call("toBe", value)
}

func Describe(name string, fn func()) {
	js.Global.Call("describe", name, fn)
}

func It(behavior string, fn func()) {
	js.Global.Call("it", behavior, fn)
}

func BeforeEach(fn func()) {
	js.Global.Call("beforeEach", fn)
}

func AfterEach(fn func()) {
	js.Global.Call("afterEach", fn)
}

func Expect(value interface{}) *Expectation {
	return &Expectation{js.Global.Call("expect", value)}
}
