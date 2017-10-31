package hamgo

import (
	"fmt"
	"net/http"
)

const (
	//DefaultPort : default port to listen
	defaultPort = "8080"
	confPort    = "port"
)

//Server : web server interface
type Server interface {
	//base
	RunAt(port string) error
	Run() error
	GetPort() string
	GetMux() *http.ServeMux
	//method
	Get(path string, handler func(ctx *WebContext)) Server
	Post(path string, handler func(ctx *WebContext)) Server
	Put(path string, handler func(ctx *WebContext)) Server
	Delete(path string, handler func(ctx *WebContext)) Server
	Head(path string, handler func(ctx *WebContext)) Server
	//get AOP
	GetBefore(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext)) Server
	GetAfter(path string, handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server
	GetBeforeAfter(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server
	//post AOP
	PostBefore(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext)) Server
	PostAfter(path string, handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server
	PostBeforeAfter(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server
	//put AOP
	PutBefore(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext)) Server
	PutAfter(path string, handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server
	PutBeforeAfter(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server
	//delete AOP
	DeleteBefore(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext)) Server
	DeleteAfter(path string, handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server
	DeleteBeforeAfter(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server
	//head AOP
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
	//filter
	Filter(handler func(ctx *WebContext) bool) *filter
}

//webServer : a web server implements Server interface
type webServer struct {
	port   string
	mux    *http.ServeMux
	filter *filter
}

//NewServer : creat a web server
func newServer() Server {
	return &webServer{mux: http.NewServeMux()}
}

//RunAt : let server run at port
func (s *webServer) RunAt(port string) error {
	s.port = ":" + port
	fmt.Printf("Server started on port :%s\n", port)
	return http.ListenAndServe(s.port, s.mux)
}

//Run : server run at default port 8080
func (s *webServer) Run() error {
	return s.RunAt(Conf.DefaultString(confPort, defaultPort))
}

//GetPort : get server run port
func (s *webServer) GetPort() string {
	return s.port
}

//
func (s *webServer) GetMux() *http.ServeMux {
	return s.mux
}

//Get : set GET method handler
func (s *webServer) Get(path string, handler func(ctx *WebContext)) Server {
	return s.Handler(path, handler, http.MethodGet)
}

//Post : set POST method handler
func (s *webServer) Post(path string, handler func(ctx *WebContext)) Server {
	return s.Handler(path, handler, http.MethodPost)
}

//Put : set PUT method handler
func (s *webServer) Put(path string, handler func(ctx *WebContext)) Server {
	return s.Handler(path, handler, http.MethodPut)
}

//Delete : set DELETE method handler
func (s *webServer) Delete(path string, handler func(ctx *WebContext)) Server {
	return s.Handler(path, handler, http.MethodDelete)
}

//Head : set HEAD method handler
func (s *webServer) Head(path string, handler func(ctx *WebContext)) Server {
	return s.Handler(path, handler, http.MethodHead)
}

//GetBefore : set func before GET method handler
func (s *webServer) GetBefore(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext)) Server {
	return s.HandlerBefore(path, handlerBefore, handler, http.MethodGet)
}

//GetAfter : set func after GET method handler
func (s *webServer) GetAfter(path string, handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server {
	return s.HandlerAfter(path, handler, handlerAfter, http.MethodGet)
}

//GetBeforeAfter : set func after & before GET method handler
func (s *webServer) GetBeforeAfter(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server {
	return s.HandlerBeforeAfter(path, handlerBefore, handler, handlerAfter, http.MethodGet)
}

//PostBefore : set func before POST method handler
func (s *webServer) PostBefore(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext)) Server {
	return s.HandlerBefore(path, handlerBefore, handler, http.MethodPost)
}

//PostAfter : set func after POST method handler
func (s *webServer) PostAfter(path string, handler func(ctx *WebContext), handlerAfter func(ctx *WebContext)) Server {
	return s.HandlerAfter(path, handler, handlerAfter, http.MethodPost)
}

//PostBeforeAfter : set func after & before POST method handler
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
	s.mux.Handle("/"+folder+"/", http.StripPrefix("/"+folder+"/", http.FileServer(http.Dir(folder))))
	return s
}

//Handler :
func (s *webServer) Handler(path string, handler func(ctx *WebContext), method string) Server {

	r := newRoute(path, method, s.filter, handler)
	s.mux.Handle(newPath(path).Route(), r)
	return s
}

//HandlerBefore :
func (s *webServer) HandlerBefore(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext), method string) Server {

	r := newBeforeRoute(path, method, s.filter, handlerBefore, handler)
	s.mux.Handle(newPath(path).Route(), r)
	return s
}

//HandlerAfter :
func (s *webServer) HandlerAfter(path string, handler func(ctx *WebContext), handlerAfter func(ctx *WebContext), method string) Server {

	r := newAfterRoute(path, method, s.filter, handler, handlerAfter)
	s.mux.Handle(newPath(path).Route(), r)
	return s
}

//HandlerBeforeAfter :
func (s *webServer) HandlerBeforeAfter(path string, handlerBefore func(ctx *WebContext), handler func(ctx *WebContext), handlerAfter func(ctx *WebContext), method string) Server {

	r := newBeforeAfterRoute(path, method, s.filter, handlerBefore, handler, handlerAfter)
	s.mux.Handle(newPath(path).Route(), r)
	return s
}

//Filter : true is pass filter , false is not pass
func (s *webServer) Filter(handler func(ctx *WebContext) bool) *filter {
	s.filter = &filter{annoURL: make([]string, 1), handler: handler}
	return s.filter
}
