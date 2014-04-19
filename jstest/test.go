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

type TestStruct struct {
	Ok int `json:"ok"`
}

func inject(deps ...interface{}) {
	js.Global.Call("inject", deps)
}

func main() {
	Describe("Test suite", func() {
		app := ng.NewModule("publicApp", []string{"ngMock", "ngRoute"})
		app.NewController("MainCtrl", func(scope *ng.Scope) {})
		app.NewController("HttpTestCtrl", func(scope *ng.Scope, http *ng.HttpService) {
			http.Get("/test").Success(func(t TestStruct, status int) {
				scope.Set("ok", t.Ok)
			})
		})
		app.Config(func(r *ng.RouteProvider) {
			r.When("/", ng.RouteOptions(
				ng.RouteController{"MainCtrl"},
				ng.RouteTemplate{"test"},
			))
		})

		BeforeEach(func() {
			js.Global.Call("module", "publicApp")
		})

		It("Should set the routes right", func() {
			ng.Bootstrap([]string{"publicApp"})
			inject("$route", "$location",
				func(route js.Object, location js.Object) {
					location.Call("path", "/")
					Expect(route.Get("routes").Get("/").Get("controller")).ToBe("MainCtrl")
				})
		})

		It("Should call the http service and get the right return", func() {
			var httpBackend, rootScope, controller js.Object
			var createController func() js.Object
			inject("$injector", func(inj js.Object) {
				httpBackend = inj.Call("get", "$httpBackend")
				httpBackend.Call("when", "GET", "/test").Call("respond",
					"{\"ok\": 1}")
				rootScope = inj.Call("get", "$rootScope")
				controller = inj.Call("get", "$controller")
				createController = func() js.Object {
					return controller.Invoke("HttpTestCtrl", map[string]js.Object{
						"$scope": rootScope,
					})
				}
			})

			ctl := createController()
			httpBackend.Call("expectGET", "/test")
			httpBackend.Call("flush")
			Expect(rootScope.Get("ok")).ToBe(1)
			_ = ctl
		})
	})
}
