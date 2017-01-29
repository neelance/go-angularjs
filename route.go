package angularjs

import (
	"reflect"

	"github.com/gopherjs/gopherjs/js"
)

func init() {
	RegisterResource(reflect.TypeOf(&RouteProvider{}), "$routeProvider", func(obj *js.Object) reflect.Value {
		return reflect.ValueOf(&RouteProvider{Object: obj})
	})

	RegisterResource(reflect.TypeOf(&RouteParams{}), "$routeParams", func(obj *js.Object) reflect.Value {
		return reflect.ValueOf(&RouteParams{Object: obj})
	})
}

// RouteProvider implements $routeProvider
type RouteProvider struct {
	*js.Object
}

// When is used to register a new URL
func (r *RouteProvider) When(url string, config RouteConfig) {
	r.Call("when", url, configToMap(config))
}

// Otherwise is used as a 404
func (r *RouteProvider) Otherwise(config RouteConfig) {
	r.Call("otherwise", configToMap(config))
}

// configToMap transforms the struct to a lowercased map
func configToMap(config RouteConfig) map[string]interface{} {
	args := map[string]interface{}{
		"controller": config.Controller,
	}

	if config.TemplateURL != "" {
		args["templateUrl"] = config.TemplateURL
	} else {
		args["template"] = config.Template
	}

	if config.Resolve != nil {

		for key, val := range config.Resolve {
			if reflect.ValueOf(val).Kind() != reflect.Func {
				panic("only funcs can be passed to RouteConfig->Resolve")
			}

			transformedFunc, err := MakeFuncInjectable(val)
			if err != nil {
				panic(err)
			}

			config.Resolve[key] = transformedFunc
		}

		args["resolve"] = config.Resolve
	}

	return args
}

// RouteConfig is an implementation of the config used to $routeProvider.when
type RouteConfig struct {
	TemplateURL string
	Template    string
	Controller  string
	Resolve     map[string]interface{}
}

// RouteParams is a copy of $routeParams
type RouteParams struct {
	*js.Object
}

// Get returns the value at the given key
func (r RouteParams) Get(key string) string {
	if r.Object.Get(key) == js.Undefined {
		return ""
	}
	return r.Object.Get(key).String()
}

// Route is the angular $route
type Route struct {
	*js.Object
}

// Params returns the $route.params value
func (r *Route) Params() map[string]string {
	ret := make(map[string]string)
	for key, val := range r.Get("params").Interface().(map[string]interface{}) {
		ret[key] = val.(string)
	}
	return ret
}
