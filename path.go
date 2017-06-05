package hamgo

import "strings"

type Path interface {
	Route() string
	PathParam(url string) map[string]string
	Paths() []string
}

type RoutePath string

const (
	PathParamPrefix = "="
)

func newPath(p string) Path {
	return RoutePath(p)
}

func (p RoutePath) Route() string {
	path := string(p)
	root := path
	if i := strings.Index(path, PathParamPrefix); i > 1 {
		root = path[0:i]
	}
	return root
}

func (p RoutePath) PathParam(url string) map[string]string {
	paths := p.Paths()
	params := RoutePath(url).Paths()
	pathParam := make(map[string]string)
	for i, param := range paths {
		if strings.HasPrefix(param, PathParamPrefix) && len(params) > i {
			pathParam[param[1:]] = params[i]
		}
	}
	return pathParam
}

func (p RoutePath) Paths() []string {
	path := string(p)
	return strings.Split(strings.Trim(path, "/ "), "/")
}
