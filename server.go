package hamgo

import (
	"fmt"
	"net/http"
	"strings"
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
	//favicon ico
	Favicon(filePath string) Server
	//common handler
	Handler(path string, handler func(ctx Context), method string) Server
	AllHandler(path string, handler func(ctx Context)) Server
	HandleFunc(path string, handler func(w http.ResponseWriter, r *http.Request)) Server
	//set filter
	AddFilter(handler func(ctx Context) bool) Server
	AddAnnoURL(url string, methods ...string) Server
	//error handler
	HTTPErrorHandler(status int, handler func(ctx Context)) Server
}

//webServer : a web server implements Server interface
type webServer struct {
	port    string
	mux     *http.ServeMux
	filters []Filter
}

var annoURLs []annoURL

type annoURL struct {
	url    string
	method string
}

//NewServer : creat a web server
func newServer() Server {
	return &webServer{mux: http.NewServeMux()}
}

//RunAt : let server run at port
func (s *webServer) RunAt(port string) error {
	s.port = ":" + port
	fmt.Printf("\nRun at port: %s\n\n", port)
	return http.ListenAndServe(s.port, s.mux)
}

//Run : server run at default port 8080
func (s *webServer) Run() error {
	if Conf != nil {
		return s.RunAt(Conf.DefaultString(confPort, defaultPort))
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
	if !isFileExist(folder) {
		panic("static folder not exist")
	}
	s.AddAnnoURL("/" + folder + "/")
	s.mux.Handle("/"+folder+"/", http.StripPrefix("/"+folder+"/", http.FileServer(http.Dir(folder))))
	return s
}

//Favicon : set "/favicon.ico"
func (s *webServer) Favicon(filePath string) Server {
	if !isFileExist(filePath) {
		panic("favicon.ico path not exist")
	}
	faviconUrl := "/favicon.ico"
	s.AddAnnoURL(faviconUrl, http.MethodPost)
	s.Get(faviconUrl, func(ctx Context) {
		ctx.File(filePath)
	})
	return s
}

//Handler : method like "POST,GET...,DELETE"
func (s *webServer) Handler(path string, handler func(ctx Context), method string) Server {
	r := newRoute(path, method, s.filters, handler)
	s.mux.Handle(newPath(path).Route(), r)
	return s
}

//HandlerAll : handel all method
func (s *webServer) AllHandler(path string, handler func(ctx Context)) Server {
	method := fmt.Sprintf("%s,%s,%s,%s,%s", http.MethodGet, http.MethodPost, http.MethodHead, http.MethodDelete, http.MethodPut)
	return s.Handler(path, handler, method)
}

//HandleFunc : handel normal func
func (s *webServer) HandleFunc(path string, handler func(w http.ResponseWriter, r *http.Request)) Server {
	s.GetMux().HandleFunc(path, handler)
	return s
}

//AddFilter : add a filter , true is pass filter , false is not pass
func (s *webServer) AddFilter(handler func(ctx Context) bool) Server {
	s.filters = append(s.filters, newFilter(handler))
	return s
}

//AddAnnoURL : add a url can pass filter
func (s *webServer) AddAnnoURL(url string, methods ...string) Server {
	if methods == nil || len(methods) == 0 {
		methods = []string{"POST", "GET", "HEAD", "DELETE", "PUT"}
	}
	annoURLs = append(annoURLs, annoURL{url, strings.ToUpper(strings.Join(methods, ","))})
	return s
}

//HTTPErrorHandler :
func (s *webServer) HTTPErrorHandler(status int, handler func(ctx Context)) Server {
	return s
}

//IsAnnoURL : check if a url is AnnoURL
func isAnnoURL(path, method string) bool {
	for _, annoURL := range annoURLs {
		if strings.HasPrefix(path, annoURL.url) && strings.Contains(annoURL.method, strings.ToUpper(method)) {
			//hase anno , pass it ,return true
			return true
		}
	}
	return false
}
