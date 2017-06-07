package main

import (
	"fmt"
	"hamgo"
	"hamgo/example/controller"
)

func main() {
	hamgo.UseConfig("./app.conf")
	hamgo.UseSession(hamgo.Hour)
	server := hamgo.New()
	server.Static("public")
	server.Get("/index/=model/=id", controller.Index)
	server.Get("/index/hello/=model/=id", controller.Hello)
	server.GetBefore("/json", controller.BeforeIndex, controller.Index)
	server.Get("/page", controller.Page)
	server.Get("/session", controller.Session)
	server.Post("/", controller.Index)
	server.Get("/bind", controller.Bind)
	fmt.Println("run at :" + hamgo.AppConfig.String("port"))
	server.RunAt(hamgo.AppConfig.String("port"))

}
