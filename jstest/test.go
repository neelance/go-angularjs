/* This package tests
through js directly with karma to verify the interactions
with Angularjs
*/
package main

import (
	ng "github.com/gopherjs/go-angularjs"
	. "github.com/gopherjs/go-angularjs/jasmine"
	"github.com/gopherjs/gopherjs/js"
)

func main() {
	Describe("Test suite", func() {
		app := ng.NewModule("publicApp", []string{"ngMock", "ngRoute"})
		app.NewController("MainCtrl", func(scope *ng.Scope) {})
		app.Config(func(services *ng.Injector) {
			ng.RouteProvider(services).When("/", ng.Options(
				ng.Route.Controller("MainCtrl"),
				ng.Route.TemplatePath("test"),
			))
		}, ng.Inject("$routeProvider"))

		BeforeEach(func() {
			js.Global.Call("module", "publicApp")
		})
		js.Global.Call("it", "Route test", func() {
			ng.Bootstrap([]string{"publicApp"})
			js.Global.Call("inject", []interface{}{"$route", "$location",
				func(route js.Object, location js.Object) {
					location.Call("path", "/")
					Expect(route.Get("routes").Get("/").Get("templateUrl")).ToBe("test")
				}})
		})
	})
}
