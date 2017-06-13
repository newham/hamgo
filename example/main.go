package main

import (
	"fmt"
	"hamgo"
	"hamgo/example/controller"
)

func main() {
	// server :=hamgo.New().UseConfig("./app.conf").UseSession(hamgo.Hour).UseLogger("./lo/app.log").Server()
	// server := hamgo.New().UseConfig("./app.conf").UseSessionByConf().UseLoggerByConf().Server()
	server := hamgo.NewUseConf("./app.conf").UseSessionByConf().UseLoggerByConf().Server()
	server.Static("public")
	server.Get("/index/=model/=id", controller.Index)
	server.Get("/index/hello/=model/=id", controller.Hello)
	server.GetBefore("/before", controller.BeforeIndex, controller.Index)
	server.Get("/json", controller.Json)
	server.Get("/page", controller.Page)
	server.Get("/session", controller.Session)
	server.Post("/", controller.Index)
	server.Get("/bind", controller.Bind)
	fmt.Println("run at :" + hamgo.Conf.String("port"))
	// server.RunAt(hamgo.Conf.String("port"))
	server.Run()

}
