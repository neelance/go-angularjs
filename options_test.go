package angularjs

import (
	"testing"
	. "launchpad.net/gocheck"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestRouteOptions(c *C) {
	str1, str2 := "mainCtrl", "html.html"
	routeOpts := Options(Route.Controller(str1),
		Route.TemplatePath(str2))
	c.Check(routeOpts.get("controller").(string), Equals, str1)
	c.Check(routeOpts.get("templateUrl").(string), Equals, str2)
}
