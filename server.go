package hamgo

import "net/http"

const (
	//DefaultPort : default port to listen
	defaultPort = "8080"
	confPort    = "port"
)

//Server :
type Server interface {
	//base
	RunAt(port string) Server
	Run() Server
	GetPort() string
	//method
	Get(path string, handler func(ctx *WebContext)) Server
	Post(path string, handler func(ctx *WebContext)) Server
	Put(path string, handler func(ctx *WebContext)) Server
	Delete(path string, handler func(ctx *WebContext)) Server
	Head(path string, handler func(ctx *WebContext)) Server
	//get inject
	GetBefore(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext)) Server
	GetAfter(path string, handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server
	GetBeforeAfter(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server
	//post inject
	PostBefore(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext)) Server
	PostAfter(path string, handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server
	PostBeforeAfter(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server
	//put inject
	PutBefore(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext)) Server
	PutAfter(path string, handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server
	PutBeforeAfter(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server
	//delete inject
	DeleteBefore(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext)) Server
	DeleteAfter(path string, handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server
	DeleteBeforeAfter(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server
	//head inject
	HeadBefore(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext)) Server
	HeadAfter(path string, handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server
	HeadBeforeAfter(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server
	//static folder
	Static(folder string) Server
	//common handler
	Handler(path string, handler func(ctx *WebContext), method string) Server
	HandlerBefore(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext), method string) Server
	HandlerAfter(path string, handler func(ctx *WebContext), handlerAfter func(ctx *WebContext), method string) Server
	HandlerBeforeAfter(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext), handlerAfter func(ctx *WebContext), method string) Server
}

//webServer :
type webServer struct {
	Port string
	Mux  *http.ServeMux
}

//NewServer : creat a web server
func newServer() Server {
	return &webServer{Mux: http.NewServeMux()}
}

//RunAt : let server run at port
func (s *webServer) RunAt(port string) Server {
	s.Port = ":" + port
	http.ListenAndServe(s.Port, s.Mux)
	return s
}

//Run :
func (s *webServer) Run() Server {
	s.Port = ":" + Conf.DefaultString(confPort, defaultPort)
	http.ListenAndServe(s.Port, s.Mux)
	return s
}

//GetPort :
func (s *webServer) GetPort() string {
	return s.Port
}

//Get :
func (s *webServer) Get(path string, handler func(ctx *WebContext)) Server {
	return s.Handler(path, handler, http.MethodGet)
}

//Post :
func (s *webServer) Post(path string, handler func(ctx *WebContext)) Server {
	return s.Handler(path, handler, http.MethodPost)
}

//Put :
func (s *webServer) Put(path string, handler func(ctx *WebContext)) Server {
	return s.Handler(path, handler, http.MethodPut)
}

//Delete :
func (s *webServer) Delete(path string, handler func(ctx *WebContext)) Server {
	return s.Handler(path, handler, http.MethodDelete)
}

//Head :
func (s *webServer) Head(path string, handler func(ctx *WebContext)) Server {
	return s.Handler(path, handler, http.MethodHead)
}

//GetBefore :
func (s *webServer) GetBefore(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext)) Server {
	return s.HandlerBefore(path, handlerBefore, handler, http.MethodGet)
}

//GetAfter :
func (s *webServer) GetAfter(path string, handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server {
	return s.HandlerAfter(path, handler, handlerAfter, http.MethodGet)
}

//GetBeforeAfter :
func (s *webServer) GetBeforeAfter(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server {
	return s.HandlerBeforeAfter(path, handlerBefore, handler, handlerAfter, http.MethodGet)
}

//PostBefore :
func (s *webServer) PostBefore(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext)) Server {
	return s.HandlerBefore(path, handlerBefore, handler, http.MethodPost)
}

//PostAfter :
func (s *webServer) PostAfter(path string, handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server {
	return s.HandlerAfter(path, handler, handlerAfter, http.MethodPost)
}

//PostBeforeAfter :
func (s *webServer) PostBeforeAfter(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server {
	return s.HandlerBeforeAfter(path, handlerBefore, handler, handlerAfter, http.MethodPost)
}

//PutBefore :
func (s *webServer) PutBefore(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext)) Server {
	return s.HandlerBefore(path, handlerBefore, handler, http.MethodGet)
}

//PutAfter :
func (s *webServer) PutAfter(path string, handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server {
	return s.HandlerAfter(path, handler, handlerAfter, http.MethodGet)
}

//PutBeforeAfter :
func (s *webServer) PutBeforeAfter(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server {
	return s.HandlerBeforeAfter(path, handlerBefore, handler, handlerAfter, http.MethodGet)
}

//DeleteBefore :
func (s *webServer) DeleteBefore(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext)) Server {
	return s.HandlerBefore(path, handlerBefore, handler, http.MethodGet)
}

//DeleteAfter :
func (s *webServer) DeleteAfter(path string, handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server {
	return s.HandlerAfter(path, handler, handlerAfter, http.MethodGet)
}

//DeleteBeforeAfter :
func (s *webServer) DeleteBeforeAfter(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server {
	return s.HandlerBeforeAfter(path, handlerBefore, handler, handlerAfter, http.MethodGet)
}

//HeadBefore :
func (s *webServer) HeadBefore(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext)) Server {
	return s.HandlerBefore(path, handlerBefore, handler, http.MethodGet)
}

//HeadAfter :
func (s *webServer) HeadAfter(path string, handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server {
	return s.HandlerAfter(path, handler, handlerAfter, http.MethodGet)
}

//HeadBeforeAfter :
func (s *webServer) HeadBeforeAfter(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server {
	return s.HandlerBeforeAfter(path, handlerBefore, handler, handlerAfter, http.MethodGet)
}

//Static :
func (s *webServer) Static(folder string) Server {

	s.Mux.Handle("/"+folder+"/", http.StripPrefix("/"+folder+"/", http.FileServer(http.Dir(folder))))
	return s
}

//Handler :
func (s *webServer) Handler(path string, handler func(ctx *WebContext), method string) Server {

	r := newRoute(path, method, handler)
	s.Mux.Handle(newPath(path).Route(), r)
	return s
}

//HandlerBefore :
func (s *webServer) HandlerBefore(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext), method string) Server {

	r := newBeforeRoute(path, method, handlerBefore, handler)
	s.Mux.Handle(newPath(path).Route(), r)
	return s
}

//HandlerAfter :
func (s *webServer) HandlerAfter(path string, handler func(ctx *WebContext), handlerAfter func(ctx *WebContext), method string) Server {

	r := newAfterRoute(path, method, handler, handlerAfter)
	s.Mux.Handle(newPath(path).Route(), r)
	return s
}

//HandlerBeforeAfter :
func (s *webServer) HandlerBeforeAfter(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext), handlerAfter func(ctx *WebContext), method string) Server {

	r := newBeforeAfterRoute(path, method, handlerBefore, handler, handlerAfter)
	s.Mux.Handle(newPath(path).Route(), r)
	return s
}
