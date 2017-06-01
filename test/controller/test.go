package controller

import "hamgo"

func Index(ctx *hamgo.Context) {
	ctx.WriteString("Index")
	ctx.Text(200)
}

func Json(ctx *hamgo.Context) {
	ctx.WriteString("Json")
	ctx.Json(200)
}
func BeforeIndex(ctx *hamgo.Context) {
	ctx.WriteString("Before ")
}

func Page(ctx *hamgo.Context) {
	ctx.WriteString("Index")
	ctx.Text(200)
}
