package hamgo

import (
	"fmt"
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
	Filter     *filter
	Func       func(ctx *WebContext)
	FuncBefore func(ctx *WebContext)
	FuncAfter  func(ctx *WebContext)
	PathKey    []string
}

type filter struct {
	handler func(ctx *WebContext) bool
	annoURL []string
}

func (f *filter) AddAnnoURL(url string) *filter {
	if len(f.annoURL) <= 1 {
		f.annoURL[0] = url
	} else {
		f.annoURL = append(f.annoURL, url)
	}
	return f
}

func doBaseFilter(route *route, r *http.Request) int {
	if !strings.Contains(route.Method, r.Method) {
		return filterMethodError
	}
	return filterOk
}

func (route *route) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch doBaseFilter(route, r) {
	case filterMethodError:
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write([]byte("405 method not allowed"))
		return
	case filter404Error:
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("404 page not found"))
		return
	}

	//set R,W
	route.W = rw
	route.R = r

	//instance context
	ctx := newWebContext(rw, r, route.Path)
	//do filter
	isFilter := true
	if route.Filter.handler != nil {
		for _, url := range route.Filter.annoURL {
			if strings.HasPrefix(route.Path, url) {
				isFilter = false
				break
			}
		}
		if isFilter && !route.Filter.handler(ctx) {
			return
		}
	}
	//do handler
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

}

func newRoute(path string, method string, filter *filter, handler func(ctx *WebContext)) *route {
	fmt.Printf("Handler : [%s]->{%s}\n", method, path)
	return &route{Inject: injectNormal, Path: path, Method: method, Filter: filter, Func: handler}
}

func newBeforeRoute(path, method string, filter *filter, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext)) *route {
	fmt.Printf("Handler : [%s]->{%s}\n", method, path)
	return &route{Inject: injectBefore, Path: path, Method: method, Filter: filter, Func: handler, FuncBefore: handlerBefore}
}

func newAfterRoute(path, method string, filter *filter, handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) *route {
	fmt.Printf("Handler : [%s]->{%s}\n", method, path)
	return &route{Inject: injectAfter, Path: path, Method: method, Filter: filter, Func: handler, FuncAfter: handlerAfter}
}

func newBeforeAfterRoute(path, method string, filter *filter, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) *route {
	fmt.Printf("Handler : [%s]->{%s}\n", method, path)
	return &route{Inject: injectBeforeAfter, Path: path, Method: method, Filter: filter, Func: handler, FuncBefore: handlerBefore, FuncAfter: handlerAfter}
}
