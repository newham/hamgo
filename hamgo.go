package hamgo

func New() IServer {
	return NewServer()
}

func SetConfig(configFile string) {
	NewConfig(configFile)
}
