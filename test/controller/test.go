package controller

import "hamgo"

func Index(ctx hamgo.IContext) {
	println("/index/")
	println("model:" + ctx.PathValue("model"))
	println("id:" + ctx.PathValue("id"))
	ctx.WriteString("Index")
	ctx.Text(200)
}

func Hello(ctx hamgo.IContext) {
	println("/index/hello/")
	println("model:" + ctx.PathValue("model"))
	println("id:" + ctx.PathValue("id"))
	ctx.WriteString("Index")
	ctx.Text(200)
}

func Json(ctx hamgo.IContext) {
	ctx.WriteString("Json")
	ctx.Json(200)
}
func BeforeIndex(ctx hamgo.IContext) {
	ctx.WriteString("Before ")
}

func Page(ctx hamgo.IContext) {
	ctx.WriteString("Index")
	ctx.Text(200)
}
