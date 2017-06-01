package main

import (
	"hamgo"
	"hamgo/test/controller"
)

func main() {
	server := hamgo.New()
	server.Static("public")
	server.Get("/index", controller.Index)
	server.GetBefore("/json", controller.BeforeIndex, controller.Index)
	server.Get("/page", controller.Page)
	server.RunAt("8087")
}
