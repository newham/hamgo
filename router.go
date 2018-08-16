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
	Inject           int
	Path             string
	Method           string
	W                http.ResponseWriter
	R                *http.Request
	Filter           Filter
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
	if route.Filter != nil {
		//not in anno url , do handler
		if !route.Filter.IsAnnoURL(route.Path) && !route.Filter.DoFilter(ctx) {
			//not pass , return false
			return false
		}
		//pass , return true
		return true
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

func newRoute(path string, method string, filter Filter, handler func(ctx Context)) *route {
	fmt.Printf("Handler : [%-4s]->{%s}\n", method, path)
	return &route{Inject: injectNormal, Path: path, Method: method, Filter: filter, Func: handler}
}
