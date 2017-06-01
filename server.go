package hamgo

import "net/http"

const (
	DEFAULT_PORT = "8080"
)

var Paths = make(map[string]string)

func New() *Server {
	return &Server{Mux: http.NewServeMux()}
}

type Server struct {
	Port string
	Mux  *http.ServeMux
	Path map[string]string
}

func (s *Server) RunAt(port string) *Server {
	s.Port = ":" + port
	http.ListenAndServe(s.Port, s.Mux)
	return s
}

func (s *Server) Run() *Server {
	s.Port = ":" + DEFAULT_PORT
	http.ListenAndServe(s.Port, s.Mux)
	return s
}

func (s *Server) Get(path string, handler func(ctx *Context)) *Server {
	return s.Handler(path, handler, http.MethodGet)
}

func (s *Server) Post(path string, handler func(ctx *Context)) *Server {
	return s.Handler(path, handler, http.MethodPost)
}

func (s *Server) Put(path string, handler func(ctx *Context)) *Server {
	return s.Handler(path, handler, http.MethodPut)
}

func (s *Server) Delete(path string, handler func(ctx *Context)) *Server {
	return s.Handler(path, handler, http.MethodDelete)
}

func (s *Server) Head(path string, handler func(ctx *Context)) *Server {
	return s.Handler(path, handler, http.MethodHead)
}

func (s *Server) GetBefore(path string, handlerBefore func(ctx *Context), handler func(ctx *Context)) *Server {
	return s.HandlerBefore(path, handlerBefore, handler, http.MethodGet)
}

func (s *Server) GetAfter(path string, handler func(ctx *Context), handlerAfter func(ctx *Context)) *Server {
	return s.HandlerAfter(path, handler, handlerAfter, http.MethodGet)
}

func (s *Server) GetBeforeAfter(path string, handlerBefore func(ctx *Context), handler func(ctx *Context), handlerAfter func(ctx *Context)) *Server {
	return s.HandlerBeforeAfter(path, handlerBefore, handler, handlerAfter, http.MethodGet)
}

//PostBefore :
func (s *Server) PostBefore(path string, handlerBefore func(ctx *Context), handler func(ctx *Context)) *Server {
	return s.HandlerBefore(path, handlerBefore, handler, http.MethodPost)
}

//PostAfter :
func (s *Server) PostAfter(path string, handler func(ctx *Context), handlerAfter func(ctx *Context)) *Server {
	return s.HandlerAfter(path, handler, handlerAfter, http.MethodPost)
}

func (s *Server) PostBeforeAfter(path string, handlerBefore func(ctx *Context), handler func(ctx *Context), handlerAfter func(ctx *Context)) *Server {
	return s.HandlerBeforeAfter(path, handlerBefore, handler, handlerAfter, http.MethodPost)
}

func (s *Server) Static(folder string) *Server {

	s.Mux.Handle("/"+folder+"/", http.StripPrefix("/"+folder+"/", http.FileServer(http.Dir(folder))))
	return s
}

func (s *Server) Handler(path string, handler func(ctx *Context), method string) *Server {
	Paths[path] = path
	r := Handler(path, method, handler)
	s.Mux.Handle(path, r)
	return s
}

func (s *Server) HandlerBefore(path string, handlerBefore func(ctx *Context), handler func(ctx *Context), method string) *Server {
	Paths[path] = path
	r := HandlerBefore(path, method, handlerBefore, handler)
	s.Mux.Handle(path, r)
	return s
}

func (s *Server) HandlerAfter(path string, handler func(ctx *Context), handlerAfter func(ctx *Context), method string) *Server {
	Paths[path] = path
	r := HandlerAfter(path, method, handler, handlerAfter)
	s.Mux.Handle(path, r)
	return s
}

func (s *Server) HandlerBeforeAfter(path string, handlerBefore func(ctx *Context), handler func(ctx *Context), handlerAfter func(ctx *Context), method string) *Server {
	Paths[path] = path
	r := HandlerBeforeAfter(path, method, handlerBefore, handler, handlerAfter)
	s.Mux.Handle(path, r)
	return s
}
