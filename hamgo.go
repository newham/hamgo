package hamgo

//Domain :
type Domain interface {
	//get server
	Server() Server
	//use modular
	UseConfig(configFile string) Domain
	UseSession(maxlifetime int64) Domain
	UseSessionByConf() Domain
	UseLogger(logFile string) Domain
	UseLoggerByConf() Domain
}

type webDomain struct {
	server Server
}

//New : create a Domain
func New() Domain {
	return &webDomain{server: newServer()}
}

//NewUseConf : create a Domain & use config
func NewUseConf(configFile string) Domain {
	d := &webDomain{server: newServer()}
	d.UseConfig(configFile)
	return d
}

func (d *webDomain) UseConfig(configFile string) Domain {
	newConfig(configFile)
	return d
}

func (d *webDomain) UseSession(maxlifetime int64) Domain {
	newSession(maxlifetime)
	return d
}

func (d *webDomain) UseSessionByConf() Domain {
	newSession(Conf.DefaultInt64(confSessionMaxTime, defaultSessionMaxTime))
	return d
}

func (d *webDomain) UseLogger(logFile string) Domain {
	newLogger(logFile)
	return d
}

func (d *webDomain) Server() Server {
	printLogo()
	return d.server
}

func (d *webDomain) UseLoggerByConf() Domain {
	newLoggerByConf()
	return d
}
