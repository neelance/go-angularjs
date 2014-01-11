package angularjs

type jsObject *struct{}

type Module struct {
	jso jsObject
}

func (m *Module) NewController(name string, constructor func(scope *Scope)) {}

const js_Module_NewController = `
	m.jso.controller(name, function($scope) {
		constructor(new Scope.Ptr($scope));
	});
`

type Scope struct {
	jso jsObject
}

func (s *Scope) Get(key string) interface{} { return nil }

const js_Scope_Get = `
	return s.jso[key];
`

func (s *Scope) GetString(key string) string { return "" }

const js_Scope_GetString = `
	return s.jso[key];
`

func (s *Scope) GetInt(key string) int { return 0 }

const js_Scope_GetInt = `
	return s.jso[key];
`

func (s *Scope) GetFloat(key string) float64 { return 0 }

const js_Scope_GetFloat = `
	return s.jso[key];
`

func (s *Scope) Set(key string, value interface{}) {}

const js_Scope_Set = `
	s.jso[key] = value;
`

func (s *Scope) Apply(f func()) {}

const js_Scope_Apply = `
	s.jso.$apply(f);
`

func (s *Scope) EvalAsync(f func()) {}

const js_Scope_EvalAsync = `
  s.jso.$evalAsync(f);
`

type JQueryElement struct {
	jso jsObject
}

func (e *JQueryElement) Prop(name string) interface{} { return nil }

const js_JQueryElement_Prop = `
	return e.jso.prop(name);
`

func (e *JQueryElement) SetProp(name, value interface{}) {}

const js_JQueryElement_SetProp = `
	e.jso.prop(name, value);
`

func (e *JQueryElement) On(events string, handler func(*Event))

const js_JQueryElement_On = `
  e.jso.on(events, function(e) { handler(new Event.Ptr(e, e.keyCode)); });
`

func (e *JQueryElement) Val() interface{} { return "" }

const js_JQueryElement_Val = `
	return e.jso.val();
`

func (e *JQueryElement) SetVal(value interface{}) {}

const js_JQueryElement_SetVal = `
	e.jso.val(value);
`

type Event struct {
	jso     jsObject
	KeyCode int
}

func (e *Event) PreventDefault() {}

const js_Event_PreventDefault = `
	e.jso.preventDefault();
`

func NewModule(name string, requires []string, configFn func()) *Module { return nil }

const js_NewModule = `
	return new Module.Ptr(angular.module(name, requires, configFn));
`

func ElementById(id string) *JQueryElement { return nil }

const js_ElementById = `
	return new JQueryElement.Ptr(angular.element(document.getElementById(id)));
`

func Service(name string) jsObject { return nil }

const js_Service = `
	return angular.element(document).injector().get(name);
`

type HttpService struct{}

var HTTP = new(HttpService)

func (s *HttpService) Get(url string, callback func(data string, status int)) {}

const js_HttpService_Get = `
	js_Service("$http").get(url).
		success(function(data, status, headers, config) {
			callback(data, status);
		}).
		error(function(data, status, headers, config) {
			callback(data, status);
		});
`
