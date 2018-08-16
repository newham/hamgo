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
	Get(path string, handler func(ctx Context)) Server
	Post(path string, handler func(ctx Context)) Server
	Put(path string, handler func(ctx Context)) Server
	Delete(path string, handler func(ctx Context)) Server
	Head(path string, handler func(ctx Context)) Server
	//static folder
	Static(folder string) Server
	//common handler
	Handler(path string, handler func(ctx Context), method string) Server
	AllHandler(path string, handler func(ctx Context)) Server
	//set filter
	Filter(handler func(ctx Context) bool) Filter
	//error handler
	HTTPErrorHandler(status int, handler func(ctx Context)) Server
}

//webServer : a web server implements Server interface
type webServer struct {
	port   string
	mux    *http.ServeMux
	filter Filter
}

//NewServer : creat a web server
func newServer() Server {
	return &webServer{mux: http.NewServeMux()}
}

//RunAt : let server run at port
func (s *webServer) RunAt(port string) error {
	s.port = ":" + port
	fmt.Printf("\nStarted at : %s\n\n", port)
	return http.ListenAndServe(s.port, s.mux)
}

//Run : server run at default port 8080
func (s *webServer) Run() error {
	if Conf != nil {
		return s.RunAt(Conf.DefaultString("port", defaultPort))
	}
	return s.RunAt(defaultPort)
}

//GetPort : get server run port
func (s *webServer) GetPort() string {
	return s.port
}

//GetMux : get http ServeMux
func (s *webServer) GetMux() *http.ServeMux {
	return s.mux
}

//Get : set GET method handler
func (s *webServer) Get(path string, handler func(ctx Context)) Server {
	return s.Handler(path, handler, http.MethodGet)
}

//Post : set POST method handler
func (s *webServer) Post(path string, handler func(ctx Context)) Server {
	return s.Handler(path, handler, http.MethodPost)
}

//Put : set PUT method handler
func (s *webServer) Put(path string, handler func(ctx Context)) Server {
	return s.Handler(path, handler, http.MethodPut)
}

//Delete : set DELETE method handler
func (s *webServer) Delete(path string, handler func(ctx Context)) Server {
	return s.Handler(path, handler, http.MethodDelete)
}

//Head : set HEAD method handler
func (s *webServer) Head(path string, handler func(ctx Context)) Server {
	return s.Handler(path, handler, http.MethodHead)
}

//Static :
func (s *webServer) Static(folder string) Server {
	s.mux.Handle("/"+folder+"/", http.StripPrefix("/"+folder+"/", http.FileServer(http.Dir(folder))))
	return s
}

//Handler :
func (s *webServer) Handler(path string, handler func(ctx Context), method string) Server {
	r := newRoute(path, method, s.filter, handler)
	s.mux.Handle(newPath(path).Route(), r)
	return s
}

//HandlerAll :
func (s *webServer) AllHandler(path string, handler func(ctx Context)) Server {
	method := fmt.Sprintf("%s,%s,%s,%s,%s", http.MethodGet, http.MethodPost, http.MethodHead, http.MethodDelete, http.MethodPut)
	return s.Handler(path, handler, method)
}

//Filter : true is pass filter , false is not pass
func (s *webServer) Filter(handler func(ctx Context) bool) Filter {
	s.filter = newFilter(handler)
	return s.filter
}

//HTTPErrorHandler :
func (s *webServer) HTTPErrorHandler(status int, handler func(ctx Context)) Server {
	return s
}
