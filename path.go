package hamgo

import "strings"

type path interface {
	Route() string
	PathParam(url string) map[string]string
	Paths() []string
}

type routePath string

const (
	pathParamPrefix = "="
)

func newPath(p string) path {
	return routePath(p)
}

func (p routePath) Route() string {
	path := string(p)
	root := path
	if i := strings.Index(path, pathParamPrefix); i > 1 {
		root = path[0:i]
	}
	return root
}

func (p routePath) PathParam(url string) map[string]string {
	paths := p.Paths()
	params := routePath(url).Paths()
	pathParam := make(map[string]string)
	for i, param := range paths {
		if strings.HasPrefix(param, pathParamPrefix) && len(params) > i {
			pathParam[param[1:]] = params[i]
		}
	}
	return pathParam
}

func (p routePath) Paths() []string {
	path := string(p)
	return strings.Split(strings.Trim(path, "/ "), "/")
}
