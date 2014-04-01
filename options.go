package angularjs

var (
	//Route is an instance of RouteSpec, it serves as a shortcut for specifying route options
	Route = &RouteSpec{NewConfigMap()}
)

//RouteSpec describes attributes of a route
type RouteSpec struct {
	*ConfigMap
}

type option func(c *ConfigMap)

type ConfigMap struct {
	opts map[string]interface{}
}

func (m *ConfigMap) Set(k string, v interface{}) {
	m.opts[k] = v
}

func (m *ConfigMap) get(k string) interface{} {
	return m.opts[k]
}

func (m *ConfigMap) Options() map[string]interface{} {
	return m.opts
}

func NewConfigMap() *ConfigMap {
	return &ConfigMap{make(map[string]interface{})}
}

//Controller sets controller option as a string for route config
func (r *RouteSpec) Controller(v string) option {
	return func(c *ConfigMap) {
		c.Set("controller", v)
	}
}

//ControllerFunc sets controller option as a func for route config
func (r *RouteSpec) ControllerFunc(v func()) option {
	return func(c *ConfigMap) {
		c.Set("controller", v)
	}
}

//TemplateFunc sets template as a function for route config
func (r *RouteSpec) TemplateFunc(v func()) option {
	return func(c *ConfigMap) {
		c.Set("template", v)
	}
}

//TemplatePath sets template as a path for route config
func (r *RouteSpec) TemplatePath(v string) option {
	return func(c *ConfigMap) {
		c.Set("templateUrl", v)
	}
}

//Options make a config from the set of options and return a ConfigMap
//representing the options
func Options(opts ...option) *ConfigMap {
	newm := NewConfigMap()
	for _, opt := range opts {
		opt(newm)
	}
	return newm
}
