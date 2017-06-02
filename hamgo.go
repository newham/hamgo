package hamgo

func New() *Server {
	return NewServer()
}

func Logger(configFile string) {
	NewConfig(configFile)
}
