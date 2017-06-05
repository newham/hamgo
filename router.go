package hamgo

import (
	"net/http"
	"strings"
)

const (
	InjectNormal      = 1
	InjectBefore      = 2
	InjectAfter       = 3
	InjectBeforeAfter = 4
	FilterMethodError = 5
	Filter404Error    = 6
	FilterOk          = 7
)

type Route struct {
	http.Handler
	Inject     int
	Path       string
	Method     string
	W          http.ResponseWriter
	R          *http.Request
	Func       func(ctx WebContext)
	FuncBefore func(ctx WebContext)
	FuncAfter  func(ctx WebContext)
	PathKey    []string
}

func filter(route *Route, r *http.Request) int {
	if !strings.Contains(route.Method, r.Method) {
		return FilterMethodError
	}
	// else if _, ok := Paths[r.URL.String()]; !ok { //not exist url
	// 	return FILTER_404_ERROR
	// }
	return FilterOk
}

func (route *Route) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch filter(route, r) {
	case FilterMethodError:
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write([]byte("405 method not allowed"))
		return
	case Filter404Error:
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("404 page not found"))
		return
	}

	//instance context
	ctx := newWebContext(rw, r, route.Path)
	switch route.Inject {
	case InjectNormal:
		route.Func(ctx)
	case InjectBefore:
		route.FuncBefore(ctx)
		route.Func(ctx)
	case InjectAfter:
		route.Func(ctx)
		route.FuncAfter(ctx)
	case InjectBeforeAfter:
		route.FuncBefore(ctx)
		route.Func(ctx)
		route.FuncAfter(ctx)
	}

	route.W = rw
	route.R = r
}

func Handler(path string, method string, handler func(ctx WebContext)) *Route {
	return &Route{Inject: InjectNormal, Path: path, Method: method, Func: handler}
}

func HandlerBefore(path, method string, handlerBefore func(ctx WebContext), handler func(ctx WebContext)) *Route {
	return &Route{Inject: InjectBefore, Path: path, Method: method, Func: handler, FuncBefore: handlerBefore}
}

func HandlerAfter(path, method string, handler func(ctx WebContext), handlerAfter func(ctx WebContext)) *Route {
	return &Route{Inject: InjectAfter, Path: path, Method: method, Func: handler, FuncAfter: handlerAfter}
}

func HandlerBeforeAfter(path, method string, handlerBefore func(ctx WebContext), handler func(ctx WebContext), handlerAfter func(ctx WebContext)) *Route {
	return &Route{Inject: InjectBeforeAfter, Path: path, Method: method, Func: handler, FuncBefore: handlerBefore, FuncAfter: handlerAfter}
}
