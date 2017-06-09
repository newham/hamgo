package hamgo

import (
	"net/http"
	"strings"
)

const (
	injectNormal      = 1
	injectBefore      = 2
	injectAfter       = 3
	injectBeforeAfter = 4
	filterMethodError = 5
	filter404Error    = 6
	filterOk          = 7
)

type route struct {
	http.Handler
	Inject     int
	Path       string
	Method     string
	W          http.ResponseWriter
	R          *http.Request
	Func       func(ctx *WebContext)
	FuncBefore func(ctx *WebContext)
	FuncAfter  func(ctx *WebContext)
	PathKey    []string
}

func filter(route *route, r *http.Request) int {
	if !strings.Contains(route.Method, r.Method) {
		return filterMethodError
	}
	return filterOk
}

func (route *route) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch filter(route, r) {
	case filterMethodError:
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write([]byte("405 method not allowed"))
		return
	case filter404Error:
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("404 page not found"))
		return
	}

	//instance context
	ctx := newWebContext(rw, r, route.Path)
	switch route.Inject {
	case injectNormal:
		route.Func(ctx)
	case injectBefore:
		route.FuncBefore(ctx)
		route.Func(ctx)
	case injectAfter:
		route.Func(ctx)
		route.FuncAfter(ctx)
	case injectBeforeAfter:
		route.FuncBefore(ctx)
		route.Func(ctx)
		route.FuncAfter(ctx)
	}

	route.W = rw
	route.R = r
}

func newRoute(path string, method string, handler func(ctx *WebContext)) *route {
	return &route{Inject: injectNormal, Path: path, Method: method, Func: handler}
}

func newBeforeRoute(path, method string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext)) *route {
	return &route{Inject: injectBefore, Path: path, Method: method, Func: handler, FuncBefore: handlerBefore}
}

func newAfterRoute(path, method string, handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) *route {
	return &route{Inject: injectAfter, Path: path, Method: method, Func: handler, FuncAfter: handlerAfter}
}

func newBeforeAfterRoute(path, method string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) *route {
	return &route{Inject: injectBeforeAfter, Path: path, Method: method, Func: handler, FuncBefore: handlerBefore, FuncAfter: handlerAfter}
}
