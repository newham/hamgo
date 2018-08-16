package main

import (
	"github.com/newham/hamgo"
	"github.com/newham/hamgo/example/controller"
)

func main() {
	// server :=hamgo.New().UseConfig("./app.conf").UseSession(hamgo.Hour).UseLogger("./lo/app.log").Server()
	// server := hamgo.New().UseConfig("./app.conf").UseSessionByConf().UseLoggerByConf().Server()
	server := hamgo.NewUseConf("./app.conf").UseSessionByConf().UseLoggerByConf().Server()
	server.Static("public")
	server.SetFilter(controller.Filter).AddAnnoURL("/login")
	server.Get("/hello", controller.Hello)
	server.Get("/login/=user/=password", controller.Login)
	server.Get("/logout", controller.Logout)
	server.Get("/index/=model/=id", controller.Index)
	server.Get("/index/hello/=model/=id", controller.Hello)
	server.Get("/json", controller.Json)
	server.Get("/page", controller.Page)
	server.Get("/session", controller.Session)
	server.Post("/", controller.Index)
	server.Get("/bind", controller.Bind)
	// server.RunAt(hamgo.Conf.String("port"))
	server.Run()

}
