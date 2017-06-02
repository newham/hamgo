package hamgo

import "strings"

type Path string

const (
	PATH_PARAM_PREFIX = "="
)

func (p Path) Root() string {
	path := string(p)
	root := path
	if i := strings.Index(path, PATH_PARAM_PREFIX); i > 1 {
		root = path[0:i]
	}
	println("root:" + root)
	return root
}

func (p Path) PathParam(url string) map[string]string {
	paths := p.Paths()
	params := Path(url).Paths()
	pathParam := make(map[string]string)
	for i, param := range paths {
		if strings.HasPrefix(param, PATH_PARAM_PREFIX) {
			pathParam[param[1:]] = params[i]
		}
	}
	return pathParam
}

func (p Path) Paths() []string {
	path := string(p)
	return strings.Split(strings.Trim(path, "/ "), "/")
}
