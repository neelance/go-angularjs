package angularjs

import (
	"encoding/json"
	"fmt"
	op "github.com/phaikawl/options"
	"reflect"

	"github.com/gopherjs/gopherjs/js"
)

type headersGetter func(string) string
type HttpTransformFunc func(string, headersGetter) string

type ReqMethod struct{ Value httpMethod }
type ReqUrl struct{ Value string }
type ReqParams struct{ Value map[string]string }
type ReqData struct{ Value string }
type ReqHeaders struct{ Value map[string]string }
type ReqTransformReqFunc struct{ Value HttpTransformFunc }
type ReqTransformRespFunc struct{ Value HttpTransformFunc }
type ReqTimeout struct{ Value int }

type ReqSpec struct {
	Method        ReqMethod            `js:"method"`
	Url           ReqUrl               `js:"url"`
	Params        ReqParams            `js:"params"`
	Data          ReqData              `js:"data"`
	Headers       ReqHeaders           `js:"headers"`
	TransformReq  ReqTransformReqFunc  `js:"transformRequest"`
	TransformResp ReqTransformRespFunc `js:"transformResponse"`
	Timeout       ReqTimeout           `js:"timeout"`
}

func jso2map(j js.Object) map[string]string {
	nm := make(map[string]string)
	if j.IsUndefined() {
		return nm
	}
	m := j.Interface().(map[string]interface{})
	for i, _ := range m {
		nm[i] = m[i].(string)
	}
	return nm
}

func ReqOptsFromJs(j js.Object) *ReqOpts {
	return HttpConfig(ReqMethod{str2httpMethod(j.Get("method").Str())},
		ReqUrl{j.Get("url").Str()},
		ReqParams{jso2map(j.Get("params"))},
		ReqData{j.Get("data").Str()},
		ReqHeaders{jso2map(j.Get("headers"))},
	)
}

type ReqOpts struct{ AngularOpts }

func (r *ReqOpts) Spec() *ReqSpec {
	return r.AngularOpts.Get().(*ReqSpec)
}

func HttpConfig(opts ...op.Option) *ReqOpts {
	return &ReqOpts{
		AngularOpts{op.NewOptions(&ReqSpec{}).Options(opts...)},
	}
}

//$http to be used at angular's run phase
type HttpService struct{ *Provider }

//$httpProvider to be used at angular's config phase
type HttpProvider struct{ *Provider }

type httpMethod string

const (
	HttpGet  httpMethod = "GET"
	HttpPost httpMethod = "POST"
)

func str2httpMethod(str string) httpMethod {
	switch str {
	case "GET":
		return HttpGet
	case "POST":
		return HttpPost
	}
	panic("Wrong method string.")
	return httpMethod("")
}

type Future struct{ js.Object }

type RequestCallback interface{}

func (ft *Future) call(state string, callback RequestCallback) *Future {
	ft.Call(state, func(data interface{}, status int, headers js.Object, config js.Object) {
		cbt := reflect.TypeOf(callback)
		in := make([]reflect.Value, cbt.NumIn())
		dparam := cbt.In(0)
		var d reflect.Value
		switch dparam.Name() {
		case "string":
			if _, ok := data.(string); !ok {
				panic("Type mismatch.")
			}
			d = reflect.ValueOf(data)
		default:
			var sdata string
			var ok bool
			if sdata, ok = data.(string); !ok {
				panic("Something is wrong.")
			}
			d = reflect.New(dparam)
			err := json.Unmarshal([]byte(sdata), d.Interface())
			if err != nil {
				panic(fmt.Sprintf("Response \"%v\" cannot be parsed to type %s. Error %v", sdata, dparam, err.Error()))
			}
		}
		in[0] = d.Elem()
		in[1] = reflect.ValueOf(status)
		reflect.ValueOf(callback).Call(in)
	})
	return ft
}

func (ft *Future) Success(callback RequestCallback) *Future {
	return ft.call("success", callback)
}
func (ft *Future) Error(callback RequestCallback) *Future {
	return ft.call("error", callback)
}

//Req performs a http request
func (s *HttpService) SimpleRequest(method httpMethod, url string) *Future {
	return s.Request(HttpConfig(ReqMethod{method}, ReqUrl{url}))
}

func (s *HttpService) Request(spec *ReqOpts) *Future {
	if !spec.IsSet("TransformReq") {
		spec.Options(ReqTransformRespFunc{func(data string, hg headersGetter) string {
			return data
		}})
	}
	future := s.Invoke(spec.toJs())
	return &Future{future}
}

func (s *HttpService) Get(url string) *Future {
	return s.SimpleRequest(HttpGet, url)
}

func AddHttpInterceptor(m *Module, name string, fn interface{}) {
	returnTypeMustBe("HttpInterceptor", fn)
	m.Call("factory", name, Ng.Inj._angularDeps(fn, func(v reflect.Value) reflect.Value {
		return v.MethodByName("ToJs").Call(make([]reflect.Value, 0))[0]
	}))
	m.Config(func(http *HttpProvider) {
		http.Object.Get("interceptors").Call("push", name)
	})
}

type HttpInterceptor struct {
	OnRequest       func(*ReqSpec)
	OnResponse      func(js.Object) interface{}
	OnRequestError  func(js.Object) interface{}
	OnResponseError func(js.Object) interface{}
}

func (hi HttpInterceptor) ToJs() (r map[string]interface{}) {
	r = make(map[string]interface{})
	if hi.OnRequest != nil {
		r["request"] = func(jso js.Object) interface{} {
			reqOpts := ReqOptsFromJs(jso)
			r := reqOpts.Spec()
			hi.OnRequest(r)
			return reqOpts.toJs()
		}
	}
	if hi.OnResponse != nil {
		r["response"] = hi.OnResponse
	}
	if hi.OnRequestError != nil {
		r["requestError"] = hi.OnRequestError
	}
	if hi.OnResponseError != nil {
		r["responseError"] = hi.OnResponseError
	}

	return
}
