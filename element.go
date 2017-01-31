package angularjs

import "github.com/gopherjs/gopherjs/js"

// Element is a basic jqLite angular element
type Element struct{ *js.Object }

// Prop calls the given property
func (e *Element) Prop(name string) *js.Object {
	return e.Call("prop", name)
}

// SetProp sets the given property
func (e *Element) SetProp(name, value interface{}) {
	e.Call("prop", name, value)
}

// On registers for events on the given element
func (e *Element) On(events string, handler func(*Event)) {
	e.Call("on", events, func(e *js.Object) {
		handler(&Event{Object: e})
	})
}

// Val returns val
func (e *Element) Val() *js.Object {
	return e.Call("val")
}

// SetVal sets val
func (e *Element) SetVal(value interface{}) {
	e.Call("val", value)
}
