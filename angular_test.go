package angularjs

import (
	"reflect"
	"testing"
	. "launchpad.net/gocheck"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestInjector(c *C) {
	test := false
	deps := Ng.Inj.AngularDeps(func(r *RouteProvider) {
		test = true
	})
	c.Check(len(deps), Equals, 2)
	fn := deps[len(deps)-1]
	in := make([]reflect.Value, len(deps)-1)
	for i, _ := range in {
		in[i] = reflect.ValueOf(&DummyJsObj{})
	}
	reflect.ValueOf(fn).Call(in)
	c.Check(test, Equals, true)
}
