package angularjs

type jsObject *struct{}

type Module struct {
	native jsObject
	SCE    *SCE
}

func (m *Module) NewController(name string, constructor func(scope *Scope)) {}

const js_Module_NewController = `
	m.native.controller(name, function($scope, $sce) {
		constructor(new Scope.Ptr($scope, new SCE.Ptr($sce)));
	});
`

type SCE struct {
	native jsObject
}

type Scope struct {
	native jsObject
}

func (s *Scope) Get(key string) interface{} { return nil }

const js_Scope_Get = `
	return s.native[key];
`

func (s *Scope) GetString(key string) string { return "" }

const js_Scope_GetString = `
	return s.native[key];
`

func (s *Scope) GetInt(key string) int { return 0 }

const js_Scope_GetInt = `
	return s.native[key];
`

func (s *Scope) GetFloat(key string) float64 { return 0 }

const js_Scope_GetFloat = `
	return s.native[key];
`

func (s *Scope) Set(key string, value interface{}) {}

const js_Scope_Set = `
	s.native[key] = value;
`

func (s *Scope) Apply(f func()) {}

const js_Scope_Apply = `
	s.native.$apply(f);
`

func (s *Scope) EvalAsync(f func()) {}

const js_Scope_EvalAsync = `
  s.native.$evalAsync(f);
`

type JQueryElement struct {
	native jsObject
}

func (e *JQueryElement) Prop(name string) string { return "" }

const js_JQueryElement_Prop = `
	return e.native.prop(name);
`

func (e *JQueryElement) SetProp(name, value interface{}) {}

const js_JQueryElement_SetProp = `
	e.native.prop(name, value);
`

func NewModule(name string, requires []string, configFn func()) *Module { return nil }

const js_NewModule = `
	return new Module.Ptr(angular.module(name, requires, configFn));
`

func ElementById(id string) *JQueryElement { return nil }

const js_ElementById = `
	return new JQueryElement.Ptr(angular.element(document.getElementById(id)));
`
