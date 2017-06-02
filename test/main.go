package main

import (
	"hamgo"
	"hamgo/test/controller"
)

func main() {
	server := hamgo.New()
	server.Static("public")
	server.Get("/index/=model/=id", controller.Index)
	server.Get("/index/hello/=model/=id", controller.Hello)
	server.GetBefore("/json", controller.BeforeIndex, controller.Index)
	server.Get("/page", controller.Page)
	server.Post("/", controller.Index)
	server.RunAt("8087")
}
