package controller

import (
	"hamgo"
	"hamgo/example/model"
	"strconv"
)

func Index(ctx *hamgo.WebContext) {
	println("/index/")
	println("model:" + ctx.PathParam("model"))
	println("id:" + ctx.PathParam("id"))
	ctx.WriteString(hamgo.Conf.String("index"))
	ctx.Text(200)
}

func Hello(ctx *hamgo.WebContext) {
	println("/index/hello/")
	println("model:" + ctx.PathParam("model"))
	println("id:" + ctx.PathParam("id"))
	ctx.WriteString("Index")
	ctx.Text(200)
}

func Json(ctx *hamgo.WebContext) {
	ctx.WriteString("Json")
	ctx.JSON(200)
}
func BeforeIndex(ctx *hamgo.WebContext) {
	ctx.WriteString("Before ")
}

func Page(ctx *hamgo.WebContext) {
	ctx.HTML(nil, "view/index.html", "view/title.tmpl")
}

func Session(ctx *hamgo.WebContext) {
	sess := ctx.GetSession()

	ct := sess.Get("countnum")
	if ct == nil {
		sess.Set("countnum", 1)
	} else {
		sess.Set("countnum", (ct.(int) + 1))
	}
	ctx.WriteString(strconv.Itoa(sess.Get("countnum").(int)))
	ctx.Text(200)
}

func Bind(ctx *hamgo.WebContext) {
	user := model.User{}
	err := ctx.BindForm(&user)

	// println("done UserName:" + user.UserName)
	hamgo.Log.Debug("done old UserName:%s", user.UserName)
	hamgo.Log.Warn("UserPassword:%s", user.UserPassword)
	hamgo.Log.Info("Age:%d", user.Age)
	hamgo.Log.Error("Email:%s", user.Email)

	if err != nil {
		hamgo.Log.Error(err.Error())
		// println(err.Error())
	}
}
