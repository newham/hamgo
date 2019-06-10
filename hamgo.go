package hamgo

//New : create a Server by Properties
func New(properties Properties) Server {
	//set server by Properties
	setLog(properties.LogFile)
	setSession(properties.SessionMaxLifeTime)
	//set logo
	printLogo()
	//return
	return newServer()
}

//New : create a Server by config
func NewByConf(configFile string) Server {
	//set server by config
	setConfig(configFile)
	return New(Properties{"", 0})
}

type Properties struct {
	LogFile            string
	SessionMaxLifeTime int64
}
