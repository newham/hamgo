package hamgo

func New() Server {
	return newServer()
}

func UseConfig(configFile string) {
	newConfig(configFile)
}

func UseSession(maxlifetime int64) {
	newSession(maxlifetime)
}

func UseLogger(logFile string) {
	newLogger(logFile)
}
