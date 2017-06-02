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
	Path       string
	Method     string
	W          http.ResponseWriter
	R          *http.Request
	Func       func(ctx IContext)
	FuncBefore func(ctx IContext)
	FuncAfter  func(ctx IContext)
	PathKey    []string
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
		rw.Write([]byte("405 method not allowed"))
		return
	case FILTER_404_ERROR:
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("404 page not found"))
		return
	}

	//instance context
	ctx := NewContext(rw, r, route.Path)
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

	route.W = rw
	route.R = r
}

func Handler(path string, method string, handler func(ctx IContext)) *Route {
	return &Route{Inject: INJECT_NORMAL, Path: path, Method: method, Func: handler}
}

func HandlerBefore(path, method string, handlerBefore func(ctx IContext), handler func(ctx IContext)) *Route {
	return &Route{Inject: INJECT_BEFORE, Path: path, Method: method, Func: handler, FuncBefore: handlerBefore}
}

func HandlerAfter(path, method string, handler func(ctx IContext), handlerAfter func(ctx IContext)) *Route {
	return &Route{Inject: INJECT_AFTER, Path: path, Method: method, Func: handler, FuncAfter: handlerAfter}
}

func HandlerBeforeAfter(path, method string, handlerBefore func(ctx IContext), handler func(ctx IContext), handlerAfter func(ctx IContext)) *Route {
	return &Route{Inject: INJECT_BEFORE_AFTER, Path: path, Method: method, Func: handler, FuncBefore: handlerBefore, FuncAfter: handlerAfter}
}
