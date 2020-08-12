package controller

import (
	"fmt"
	"strconv"

	"github.com/newham/hamgo"
	"github.com/newham/hamgo/example/model"
)

const (
	USER_SESSION string = "userSession"
)

func LoginByForm(ctx hamgo.Context) {
	usr := ctx.FormValue("usr")
	pwd := ctx.FormValue("pwd")
	if usr == "admin" && pwd == "123456" {
		session := ctx.GetSession()
		session.Set(USER_SESSION, usr)
		hamgo.Log.Info("Login")
		ctx.WriteString("Login Success!")
		ctx.Redirect("/hello")
	} else {
		ctx.WriteString("Login failed!")
		ctx.Text(200)
	}
}

func Login(ctx hamgo.Context) {
	user := ctx.PathParam("user")
	password := ctx.PathParam("password")
	if user == "admin" && password == "123456" {
		ctx.GetSession().Set(USER_SESSION, user)
		if ctx.Method() == "GET" {
			ctx.Redirect("/json")
			return
		}
		ctx.WriteString("login success")
		ctx.Text(200)
		hamgo.Log.Info("%s login success", user)
	} else {
		ctx.WriteString("login failed")
		ctx.Text(400)
		hamgo.Log.Error("%s login failed", user)
	}

}

func Logout(ctx hamgo.Context) {
	hamgo.Log.Info("logout:%s", ctx.GetSession().Get(USER_SESSION))
	ctx.DeleteSession()
	ctx.WriteString("logout success")
	ctx.Text(200)
}

func Index(ctx hamgo.Context) {
	println("model:" + ctx.PathParam("model"))
	println("id:" + ctx.PathParam("id"))
	ctx.WriteString(hamgo.Conf.String("index"))
	ctx.Text(200)
}

func Filter(ctx hamgo.Context) bool {
	session := ctx.GetSession()
	if session != nil && session.Get(USER_SESSION) != nil {
		return true
	}

	ctx.PutData("code", 401)
	ctx.PutData("msg", "Unauthorized")
	// ctx.JSON(401, nil)
	ctx.Redirect("/page")
	hamgo.Log.Error("401,%s", "Unauthorized")
	return false
}

func Hello(ctx hamgo.Context) {
	ctx.WriteString("Hello World!\n")
	ctx.WriteString(fmt.Sprintf("Session Left Time:%d", ctx.GetSession().LeftTime()))
	ctx.Text(200)
}

func Refresh(ctx hamgo.Context) {
	ctx.RefreshSession()
	ctx.Redirect("/hello")
}

func Json(ctx hamgo.Context) {
	ctx.PutData("say", "hello world")
	ctx.PutData("leftTime", ctx.GetSession().LeftTime())
	ctx.JSON(200, nil)
}

func JsonFromData(ctx hamgo.Context) {
	data := map[string]interface{}{"say": "hello world"}
	ctx.JSONFrom(200, data)
}

func Page(ctx hamgo.Context) {
	ctx.PutData("Title", "Hello World")
	ctx.HTML("view/index.html", "view/title.tmpl")
}

func Session(ctx hamgo.Context) {
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

func Bind(ctx hamgo.Context) {
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
