package controller

import (
	"hamgo"
	"hamgo/example/model"
	"strconv"
)

const (
	USER_SESSION string = "userSession"
)

func Login(ctx *hamgo.WebContext) {
	println("/login/")
	user := ctx.PathParam("user")
	password := ctx.PathParam("password")
	if user == "admin" && password == "123456" {
		ctx.GetSession().Set(USER_SESSION, user)
		ctx.WriteString("login success")
		ctx.Text(200)
	} else {
		ctx.WriteString("login failed")
		ctx.Text(400)
	}

}

func Logout(ctx *hamgo.WebContext) {
	hamgo.Log.Info("logout:%s", ctx.GetSession().Get(USER_SESSION))
	ctx.GetSession().Delete(USER_SESSION)
	ctx.WriteString("logout success")
	ctx.Text(200)
}

func Index(ctx *hamgo.WebContext) {
	println("/index/")
	println("model:" + ctx.PathParam("model"))
	println("id:" + ctx.PathParam("id"))
	ctx.WriteString(hamgo.Conf.String("index"))
	ctx.Text(200)
}

func Filter(ctx *hamgo.WebContext) bool {
	if ctx.GetSession().Get(USER_SESSION) != nil {
		return true
	}
	ctx.PutData("code", 401)
	ctx.PutData("msg", "Unauthorized")
	ctx.JSON(401)
	hamgo.Log.Error("401,%s", "Unauthorized")
	return false
}

func Hello(ctx *hamgo.WebContext) {
	println("/index/hello/")
	println("model:" + ctx.PathParam("model"))
	println("id:" + ctx.PathParam("id"))
	ctx.WriteString("Index")
	ctx.Text(200)
}

func Json(ctx *hamgo.WebContext) {
	ctx.PutData("say", "hello world")
	ctx.JSON(200)
}

func JsonFromData(ctx *hamgo.WebContext) {
	data := map[string]interface{}{"say": "hello world"}
	ctx.JSONFrom(200, data)
}

func BeforeIndex(ctx *hamgo.WebContext) {
	ctx.WriteString("Before ")
}

func Page(ctx *hamgo.WebContext) {
	ctx.PutData("Title", "Hello World")
	ctx.HTML("view/index.html", "view/title.tmpl")
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
	errs := ctx.BindForm(&user)

	// println("done UserName:" + user.UserName)
	hamgo.Log.Debug("done old UserName:%s", user.UserName)
	hamgo.Log.Warn("UserPassword:%s", user.UserPassword)
	hamgo.Log.Info("Age:%d", user.Age)
	hamgo.Log.Error("Email:%s", user.Email)
	//
	hamgo.Log.Error("bindForm result--------")
	if len(errs) > 0 {
		for k, err := range errs {
			hamgo.Log.Error("%s:%s", k, err.Error())
		}
	}
}
