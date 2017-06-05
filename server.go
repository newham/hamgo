package hamgo

import "net/http"

const (
	//DefaultPort : default port to listen
	DefaultPort = "8080"
)

//Server :
type Server interface {
	//base
	RunAt(port string) *WebServer
	Run() *WebServer
	GetPort() string
	//method
	Get(path string, handler func(ctx WebContext)) *WebServer
	Post(path string, handler func(ctx WebContext)) *WebServer
	Put(path string, handler func(ctx WebContext)) *WebServer
	Delete(path string, handler func(ctx WebContext)) *WebServer
	Head(path string, handler func(ctx WebContext)) *WebServer
	//get inject
	GetBefore(path string, handlerBefore func(ctx WebContext), handler func(ctx WebContext)) *WebServer
	GetAfter(path string, handler func(ctx WebContext), handlerAfter func(ctx WebContext)) *WebServer
	GetBeforeAfter(path string, handlerBefore func(ctx WebContext), handler func(ctx WebContext), handlerAfter func(ctx WebContext)) *WebServer
	//post inject
	PostBefore(path string, handlerBefore func(ctx WebContext), handler func(ctx WebContext)) *WebServer
	PostAfter(path string, handler func(ctx WebContext), handlerAfter func(ctx WebContext)) *WebServer
	PostBeforeAfter(path string, handlerBefore func(ctx WebContext), handler func(ctx WebContext), handlerAfter func(ctx WebContext)) *WebServer
	//put inject
	PutBefore(path string, handlerBefore func(ctx WebContext), handler func(ctx WebContext)) *WebServer
	PutAfter(path string, handler func(ctx WebContext), handlerAfter func(ctx WebContext)) *WebServer
	PutBeforeAfter(path string, handlerBefore func(ctx WebContext), handler func(ctx WebContext), handlerAfter func(ctx WebContext)) *WebServer
	//delete inject
	DeleteBefore(path string, handlerBefore func(ctx WebContext), handler func(ctx WebContext)) *WebServer
	DeleteAfter(path string, handler func(ctx WebContext), handlerAfter func(ctx WebContext)) *WebServer
	DeleteBeforeAfter(path string, handlerBefore func(ctx WebContext), handler func(ctx WebContext), handlerAfter func(ctx WebContext)) *WebServer
	//head inject
	HeadBefore(path string, handlerBefore func(ctx WebContext), handler func(ctx WebContext)) *WebServer
	HeadAfter(path string, handler func(ctx WebContext), handlerAfter func(ctx WebContext)) *WebServer
	HeadBeforeAfter(path string, handlerBefore func(ctx WebContext), handler func(ctx WebContext), handlerAfter func(ctx WebContext)) *WebServer
	//static folder
	Static(folder string) *WebServer
	//common handler
	Handler(path string, handler func(ctx WebContext), method string) *WebServer
	HandlerBefore(path string, handlerBefore func(ctx WebContext), handler func(ctx WebContext), method string) *WebServer
	HandlerAfter(path string, handler func(ctx WebContext), handlerAfter func(ctx WebContext), method string) *WebServer
	HandlerBeforeAfter(path string, handlerBefore func(ctx WebContext), handler func(ctx WebContext), handlerAfter func(ctx WebContext), method string) *WebServer
}

//WebServer :
type WebServer struct {
	Port string
	Mux  *http.ServeMux
}

//NewServer : creat a web server
func newServer() Server {
	return &WebServer{Mux: http.NewServeMux()}
}

//RunAt : let server run at port
func (s *WebServer) RunAt(port string) *WebServer {
	s.Port = ":" + port
	http.ListenAndServe(s.Port, s.Mux)
	return s
}

//Run :
func (s *WebServer) Run() *WebServer {
	s.Port = ":" + DefaultPort
	http.ListenAndServe(s.Port, s.Mux)
	return s
}

//GetPort :
func (s *WebServer) GetPort() string {
	return s.Port
}

//Get :
func (s *WebServer) Get(path string, handler func(ctx WebContext)) *WebServer {
	return s.Handler(path, handler, http.MethodGet)
}

//Post :
func (s *WebServer) Post(path string, handler func(ctx WebContext)) *WebServer {
	return s.Handler(path, handler, http.MethodPost)
}

//Put :
func (s *WebServer) Put(path string, handler func(ctx WebContext)) *WebServer {
	return s.Handler(path, handler, http.MethodPut)
}

//Delete :
func (s *WebServer) Delete(path string, handler func(ctx WebContext)) *WebServer {
	return s.Handler(path, handler, http.MethodDelete)
}

//Head :
func (s *WebServer) Head(path string, handler func(ctx WebContext)) *WebServer {
	return s.Handler(path, handler, http.MethodHead)
}

//GetBefore :
func (s *WebServer) GetBefore(path string, handlerBefore func(ctx WebContext), handler func(ctx WebContext)) *WebServer {
	return s.HandlerBefore(path, handlerBefore, handler, http.MethodGet)
}

//GetAfter :
func (s *WebServer) GetAfter(path string, handler func(ctx WebContext), handlerAfter func(ctx WebContext)) *WebServer {
	return s.HandlerAfter(path, handler, handlerAfter, http.MethodGet)
}

//GetBeforeAfter :
func (s *WebServer) GetBeforeAfter(path string, handlerBefore func(ctx WebContext), handler func(ctx WebContext), handlerAfter func(ctx WebContext)) *WebServer {
	return s.HandlerBeforeAfter(path, handlerBefore, handler, handlerAfter, http.MethodGet)
}

//PostBefore :
func (s *WebServer) PostBefore(path string, handlerBefore func(ctx WebContext), handler func(ctx WebContext)) *WebServer {
	return s.HandlerBefore(path, handlerBefore, handler, http.MethodPost)
}

//PostAfter :
func (s *WebServer) PostAfter(path string, handler func(ctx WebContext), handlerAfter func(ctx WebContext)) *WebServer {
	return s.HandlerAfter(path, handler, handlerAfter, http.MethodPost)
}

//PostBeforeAfter :
func (s *WebServer) PostBeforeAfter(path string, handlerBefore func(ctx WebContext), handler func(ctx WebContext), handlerAfter func(ctx WebContext)) *WebServer {
	return s.HandlerBeforeAfter(path, handlerBefore, handler, handlerAfter, http.MethodPost)
}

//PutBefore :
func (s *WebServer) PutBefore(path string, handlerBefore func(ctx WebContext), handler func(ctx WebContext)) *WebServer {
	return s.HandlerBefore(path, handlerBefore, handler, http.MethodGet)
}

//PutAfter :
func (s *WebServer) PutAfter(path string, handler func(ctx WebContext), handlerAfter func(ctx WebContext)) *WebServer {
	return s.HandlerAfter(path, handler, handlerAfter, http.MethodGet)
}

//PutBeforeAfter :
func (s *WebServer) PutBeforeAfter(path string, handlerBefore func(ctx WebContext), handler func(ctx WebContext), handlerAfter func(ctx WebContext)) *WebServer {
	return s.HandlerBeforeAfter(path, handlerBefore, handler, handlerAfter, http.MethodGet)
}

//DeleteBefore :
func (s *WebServer) DeleteBefore(path string, handlerBefore func(ctx WebContext), handler func(ctx WebContext)) *WebServer {
	return s.HandlerBefore(path, handlerBefore, handler, http.MethodGet)
}

//DeleteAfter :
func (s *WebServer) DeleteAfter(path string, handler func(ctx WebContext), handlerAfter func(ctx WebContext)) *WebServer {
	return s.HandlerAfter(path, handler, handlerAfter, http.MethodGet)
}

//DeleteBeforeAfter :
func (s *WebServer) DeleteBeforeAfter(path string, handlerBefore func(ctx WebContext), handler func(ctx WebContext), handlerAfter func(ctx WebContext)) *WebServer {
	return s.HandlerBeforeAfter(path, handlerBefore, handler, handlerAfter, http.MethodGet)
}

//HeadBefore :
func (s *WebServer) HeadBefore(path string, handlerBefore func(ctx WebContext), handler func(ctx WebContext)) *WebServer {
	return s.HandlerBefore(path, handlerBefore, handler, http.MethodGet)
}

//HeadAfter :
func (s *WebServer) HeadAfter(path string, handler func(ctx WebContext), handlerAfter func(ctx WebContext)) *WebServer {
	return s.HandlerAfter(path, handler, handlerAfter, http.MethodGet)
}

//HeadBeforeAfter :
func (s *WebServer) HeadBeforeAfter(path string, handlerBefore func(ctx WebContext), handler func(ctx WebContext), handlerAfter func(ctx WebContext)) *WebServer {
	return s.HandlerBeforeAfter(path, handlerBefore, handler, handlerAfter, http.MethodGet)
}

//Static :
func (s *WebServer) Static(folder string) *WebServer {

	s.Mux.Handle("/"+folder+"/", http.StripPrefix("/"+folder+"/", http.FileServer(http.Dir(folder))))
	return s
}

//Handler :
func (s *WebServer) Handler(path string, handler func(ctx WebContext), method string) *WebServer {

	r := Handler(path, method, handler)
	s.Mux.Handle(newPath(path).Route(), r)
	return s
}

//HandlerBefore :
func (s *WebServer) HandlerBefore(path string, handlerBefore func(ctx WebContext), handler func(ctx WebContext), method string) *WebServer {

	r := HandlerBefore(path, method, handlerBefore, handler)
	s.Mux.Handle(newPath(path).Route(), r)
	return s
}

//HandlerAfter :
func (s *WebServer) HandlerAfter(path string, handler func(ctx WebContext), handlerAfter func(ctx WebContext), method string) *WebServer {

	r := HandlerAfter(path, method, handler, handlerAfter)
	s.Mux.Handle(newPath(path).Route(), r)
	return s
}

//HandlerBeforeAfter :
func (s *WebServer) HandlerBeforeAfter(path string, handlerBefore func(ctx WebContext), handler func(ctx WebContext), handlerAfter func(ctx WebContext), method string) *WebServer {

	r := HandlerBeforeAfter(path, method, handlerBefore, handler, handlerAfter)
	s.Mux.Handle(newPath(path).Route(), r)
	return s
}
