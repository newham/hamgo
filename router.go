package hamgo

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	filterMethodError = 5
	filter404Error    = 6
)

type route struct {
	http.Handler
	Path             string
	Method           string
	W                http.ResponseWriter
	R                *http.Request
	Filters          []Filter
	Func             func(ctx Context)
	PathKey          []string
	HTTPErrorHandler func(ctx Context)
}

func (route *route) doFilter(rw http.ResponseWriter, r *http.Request) (Context, bool) {
	//instance context
	ctx := newWebContext(rw, r, route.Path)
	//do user filter
	if !route.doUserFilter(ctx) {
		return ctx, false
	}
	//do base filter
	if !route.doBaseFilter(rw, r) {
		return ctx, false
	}
	return ctx, true
}

func (route *route) doBaseFilter(rw http.ResponseWriter, r *http.Request) bool {
	var filterError int
	if !strings.Contains(route.Method, r.Method) {
		filterError = filterMethodError
	}
	switch filterError {
	case filterMethodError:
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write([]byte("405 method not allowed"))
		return false
	case filter404Error:
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("404 page not found"))
		return false
	}
	return true
}

func (route *route) doUserFilter(ctx Context) bool {
	if isAnnoURL(route.Path, ctx.Method()) {
		return true
	}
	for _, filter := range route.Filters {
		if !filter.DoFilter(ctx) {
			//not pass , return false
			return false
		}
	}
	//not set filter , pass
	return true
}

func (route *route) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	//do filter
	ctx, filterResult := route.doFilter(rw, r)
	if !filterResult {
		return
	}
	//set R,W
	route.W = rw
	route.R = r
	//do handler
	route.Func(ctx)
}

func newRoute(path string, method string, filters []Filter, handler func(ctx Context)) *route {
	fmt.Printf("Handler: [%-4s]->{%s}\n", method, path)
	return &route{Path: path, Method: method, Filters: filters, Func: handler}
}
