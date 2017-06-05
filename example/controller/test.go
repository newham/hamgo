package controller

import (
	"hamgo"
	"hamgo/example/model"
	"strconv"
	"time"
)

func Index(ctx hamgo.WebContext) {
	println("/index/")
	println("model:" + ctx.PathParam("model"))
	println("id:" + ctx.PathParam("id"))
	ctx.WriteString(hamgo.AppConfig.String("index"))
	ctx.Text(200)
}

func Hello(ctx hamgo.WebContext) {
	println("/index/hello/")
	println("model:" + ctx.PathParam("model"))
	println("id:" + ctx.PathParam("id"))
	ctx.WriteString("Index")
	ctx.Text(200)
}

func Json(ctx hamgo.WebContext) {
	ctx.WriteString("Json")
	ctx.Json(200)
}
func BeforeIndex(ctx hamgo.WebContext) {
	ctx.WriteString("Before ")
}

func Page(ctx hamgo.WebContext) {
	ctx.Html(nil, "view/index.html", "view/title.tmpl")
}

func Session(ctx hamgo.WebContext) {
	sess := ctx.GetSession()
	createtime := sess.Get("createtime")
	if createtime == nil {
		sess.Set("createtime", time.Now().Unix())
	} else if (createtime.(int64) + 360) < (time.Now().Unix()) {
		hamgo.Sessions.SessionDestroy(ctx.W, ctx.R)
		sess = hamgo.Sessions.SessionStart(ctx.W, ctx.R)
	}
	ct := sess.Get("countnum")
	if ct == nil {
		sess.Set("countnum", 1)
	} else {
		sess.Set("countnum", (ct.(int) + 1))
	}
	ctx.WriteString(strconv.Itoa(sess.Get("countnum").(int)))
	ctx.Text(200)
}

func Bind(ctx hamgo.WebContext) {
	user := model.User{}
	u := ctx.BindForm(user).(*model.User)
	println("done UserName:" + u.UserName)
}
