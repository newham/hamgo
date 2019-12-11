package hamgo

//New : create a Server by Properties
func New(properties *Properties) Server {
	//set server by Properties
	if properties != nil {
		setLog(properties.LogFile)
		setSession(properties.SessionMaxLifeTime)
	}
	//set logo
	printLogo()
	//return
	return newServer()
}

//New : create a Server by config
func NewByConf(configFile string) Server {
	//set server by config
	setConfig(configFile)
	return New(&Properties{"", Conf.DefaultInt64("session_max_time", 0)})
}

func NewProperties(logFile string, sessionMaxLifeTime int64) *Properties {
	return &Properties{logFile, sessionMaxLifeTime}
}

type Properties struct {
	LogFile            string
	SessionMaxLifeTime int64
}
