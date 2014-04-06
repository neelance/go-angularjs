package angularjs

import op "github.com/phaikawl/options"

type RouteController struct{ Value string }
type RouteTemplate struct{ Value string }

type RouteSpec struct {
	Controller RouteController `js:"controller"`
	Template   RouteTemplate   `js:"templateUrl"`
}

type AngularOpts struct{ *op.OptionsProvider }

func (o *AngularOpts) toJs() map[string]interface{} {
	return o.OptionsProvider.ExportToMapWithTag("js")
}

func RouteOptions(opts ...op.Option) *AngularOpts {
	return &AngularOpts{op.NewOptions(&RouteSpec{}).Options(opts...)}
}
