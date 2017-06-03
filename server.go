package hamgo

import "net/http"

const (
	DEFAULT_PORT = "8080"
)

type IServer interface {
	RunAt(port string) *Server
	Run() *Server
	Get(path string, handler func(ctx IContext)) *Server
	Post(path string, handler func(ctx IContext)) *Server
	Put(path string, handler func(ctx IContext)) *Server
	Delete(path string, handler func(ctx IContext)) *Server
	Head(path string, handler func(ctx IContext)) *Server
	GetBefore(path string, handlerBefore func(ctx IContext), handler func(ctx IContext)) *Server
	GetAfter(path string, handler func(ctx IContext), handlerAfter func(ctx IContext)) *Server
	GetBeforeAfter(path string, handlerBefore func(ctx IContext), handler func(ctx IContext), handlerAfter func(ctx IContext)) *Server
	PostBefore(path string, handlerBefore func(ctx IContext), handler func(ctx IContext)) *Server
	PostAfter(path string, handler func(ctx IContext), handlerAfter func(ctx IContext)) *Server
	PostBeforeAfter(path string, handlerBefore func(ctx IContext), handler func(ctx IContext), handlerAfter func(ctx IContext)) *Server
	Static(folder string) *Server
	Handler(path string, handler func(ctx IContext), method string) *Server
	HandlerBefore(path string, handlerBefore func(ctx IContext), handler func(ctx IContext), method string) *Server
	HandlerAfter(path string, handler func(ctx IContext), handlerAfter func(ctx IContext), method string) *Server
	HandlerBeforeAfter(path string, handlerBefore func(ctx IContext), handler func(ctx IContext), handlerAfter func(ctx IContext), method string) *Server
}
type Server struct {
	Port string
	Mux  *http.ServeMux
}

func NewServer() IServer {
	return &Server{Mux: http.NewServeMux()}
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

func (s *Server) Get(path string, handler func(ctx IContext)) *Server {
	return s.Handler(path, handler, http.MethodGet)
}

func (s *Server) Post(path string, handler func(ctx IContext)) *Server {
	return s.Handler(path, handler, http.MethodPost)
}

func (s *Server) Put(path string, handler func(ctx IContext)) *Server {
	return s.Handler(path, handler, http.MethodPut)
}

func (s *Server) Delete(path string, handler func(ctx IContext)) *Server {
	return s.Handler(path, handler, http.MethodDelete)
}

func (s *Server) Head(path string, handler func(ctx IContext)) *Server {
	return s.Handler(path, handler, http.MethodHead)
}

func (s *Server) GetBefore(path string, handlerBefore func(ctx IContext), handler func(ctx IContext)) *Server {
	return s.HandlerBefore(path, handlerBefore, handler, http.MethodGet)
}

func (s *Server) GetAfter(path string, handler func(ctx IContext), handlerAfter func(ctx IContext)) *Server {
	return s.HandlerAfter(path, handler, handlerAfter, http.MethodGet)
}

func (s *Server) GetBeforeAfter(path string, handlerBefore func(ctx IContext), handler func(ctx IContext), handlerAfter func(ctx IContext)) *Server {
	return s.HandlerBeforeAfter(path, handlerBefore, handler, handlerAfter, http.MethodGet)
}

//PostBefore :
func (s *Server) PostBefore(path string, handlerBefore func(ctx IContext), handler func(ctx IContext)) *Server {
	return s.HandlerBefore(path, handlerBefore, handler, http.MethodPost)
}

//PostAfter :
func (s *Server) PostAfter(path string, handler func(ctx IContext), handlerAfter func(ctx IContext)) *Server {
	return s.HandlerAfter(path, handler, handlerAfter, http.MethodPost)
}

func (s *Server) PostBeforeAfter(path string, handlerBefore func(ctx IContext), handler func(ctx IContext), handlerAfter func(ctx IContext)) *Server {
	return s.HandlerBeforeAfter(path, handlerBefore, handler, handlerAfter, http.MethodPost)
}

func (s *Server) Static(folder string) *Server {

	s.Mux.Handle("/"+folder+"/", http.StripPrefix("/"+folder+"/", http.FileServer(http.Dir(folder))))
	return s
}

func (s *Server) Handler(path string, handler func(ctx IContext), method string) *Server {

	r := Handler(path, method, handler)
	s.Mux.Handle(Path(path).Root(), r)
	return s
}

func (s *Server) HandlerBefore(path string, handlerBefore func(ctx IContext), handler func(ctx IContext), method string) *Server {

	r := HandlerBefore(path, method, handlerBefore, handler)
	s.Mux.Handle(Path(path).Root(), r)
	return s
}

func (s *Server) HandlerAfter(path string, handler func(ctx IContext), handlerAfter func(ctx IContext), method string) *Server {

	r := HandlerAfter(path, method, handler, handlerAfter)
	s.Mux.Handle(Path(path).Root(), r)
	return s
}

func (s *Server) HandlerBeforeAfter(path string, handlerBefore func(ctx IContext), handler func(ctx IContext), handlerAfter func(ctx IContext), method string) *Server {

	r := HandlerBeforeAfter(path, method, handlerBefore, handler, handlerAfter)
	s.Mux.Handle(Path(path).Root(), r)
	return s
}
