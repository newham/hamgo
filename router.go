package hamgo

import (
	"net/http"
	"strings"
)

const (
	INJECT_NORMAL       = 1
	INJECT_BEFORE       = 2
	INJECT_AFTER        = 3
	INJECT_BEFORE_AFTER = 4
	FILTER_METHOD_ERROR = 5
	FILTER_404_ERROR    = 6
	FILTER_OK           = 7
)

type Route struct {
	http.Handler
	Inject     int
	Pattern    string
	Method     string
	Resp       http.ResponseWriter
	Req        *http.Request
	Func       func(ctx *Context)
	FuncBefore func(ctx *Context)
	FuncAfter  func(ctx *Context)
}

func Handler(pattern, method string, handler func(ctx *Context)) *Route {
	return &Route{Inject: INJECT_NORMAL, Pattern: pattern, Method: method, Func: handler}
}
func filter(route *Route, r *http.Request) int {
	if !strings.Contains(route.Method, r.Method) {
		return FILTER_METHOD_ERROR
	}
	// else if _, ok := Paths[r.URL.String()]; !ok { //not exist url
	// 	return FILTER_404_ERROR
	// }
	return FILTER_OK
}
func (route *Route) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch filter(route, r) {
	case FILTER_METHOD_ERROR:
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write([]byte("Wrong Method"))
		return
	case FILTER_404_ERROR:
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("<h1>404 Not found page</h1>"))
		return
	}

	//instance context
	ctx := NewContext(rw, r)
	switch route.Inject {
	case INJECT_NORMAL:
		route.Func(ctx)
	case INJECT_BEFORE:
		route.FuncBefore(ctx)
		route.Func(ctx)
	case INJECT_AFTER:
		route.Func(ctx)
		route.FuncAfter(ctx)
	case INJECT_BEFORE_AFTER:
		route.FuncBefore(ctx)
		route.Func(ctx)
		route.FuncAfter(ctx)
	}

	route.Resp = rw
	route.Req = r
}

func HandlerBefore(pattern, method string, handlerBefore func(ctx *Context), handler func(ctx *Context)) *Route {
	return &Route{Inject: INJECT_BEFORE, Pattern: pattern, Method: method, Func: handler, FuncBefore: handlerBefore}
}

func HandlerAfter(pattern, method string, handler func(ctx *Context), handlerAfter func(ctx *Context)) *Route {
	return &Route{Inject: INJECT_AFTER, Pattern: pattern, Method: method, Func: handler, FuncAfter: handlerAfter}
}

func HandlerBeforeAfter(pattern, method string, handlerBefore func(ctx *Context), handler func(ctx *Context), handlerAfter func(ctx *Context)) *Route {
	return &Route{Inject: INJECT_BEFORE_AFTER, Pattern: pattern, Method: method, Func: handler, FuncBefore: handlerBefore, FuncAfter: handlerAfter}
}
