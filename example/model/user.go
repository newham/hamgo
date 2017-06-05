package model

type User struct {
	UserName     string `form:"username" regexp:""`
	UserPassword string `form:"password" filter:""`
}
