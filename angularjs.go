package angularjs

func NewModule(name string, requires []string, configFn func()) *Module { return nil }

const js_NewModule = `
	requires = requires ? go$sliceToArray(requires) : [];
	return new Module.Ptr(angular.module(name, requires, configFn));
`

type Module struct {
	native interface{}
	SCE    *SCE
}

func (m *Module) NewController(name string, constructor func(scope *Scope)) {}

const js_Module_NewController = `
	this.native.controller(name, function($scope, $sce) {
		constructor(new Scope.Ptr($scope, new SCE.Ptr($sce)));
	});
`

func (m *Module) NewFilter(name string, fn func(text string, arguments []string) string) {}

const js_Module_NewFilter = `
	this.native.filter(name, function() {
		return function(text) {
			return fn(text, new (go$sliceType(Go$String))(Array.prototype.slice.call(arguments, 1)));
		};
	});
`

type Scope struct {
	native interface{}
}

func (s *Scope) GetString(key string) string { return "" }

const js_Scope_GetString = `
	return go$internalizeString(String(this.native[key]));
`

func (s *Scope) GetInt(key string) int { return 0 }

const js_Scope_GetInt = `
	return parseInt(this.native[key]);
`

func (s *Scope) GetFloat(key string) float64 { return 0 }

const js_Scope_GetFloat = `
	return parseFloat(this.native[key]);
`

func (s *Scope) Set(key string, value interface{}) {}

const js_Scope_Set = `
	if (value.constructor === Go$String) {
		this.native[key] = go$externalizeString(value.go$val);
		return;
	}
	if (value.array !== undefined) { // TODO we need a better solution here
		this.native[key] = go$sliceToArray(value);
		return;
	}
	this.native[key] = value.go$val;
`

func (s *Scope) Apply(f func()) {}

const js_Scope_Apply = `
	this.native.$apply(f);
`

type SCE struct {
	native interface{}
}
