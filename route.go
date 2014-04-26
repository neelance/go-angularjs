package angularjs

import op "github.com/phaikawl/options"

type RouteController struct{ Value string }
type RouteTemplate struct{ Value string }

type RouteSpec struct {
	Controller RouteController `js:"controller"`
	Template   RouteTemplate   `js:"templateUrl"`
}

type RouteOpts struct{ AngularOpts }

func RouteOptions(opts ...op.Option) *RouteOpts {
	return &RouteOpts{
		AngularOpts{op.NewOptions(&RouteSpec{}).Options(opts...)},
	}
}

type RouteProvider struct {
	*Provider
}

func (r *RouteProvider) When(path string, route *RouteOpts) *RouteProvider {
	r.Call("when", path, route.toJs())
	return r
}

func (r *RouteProvider) Otherwise(route *RouteOpts) {
	r.Call("otherwise", route.toJs())
}
