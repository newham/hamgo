package model

type User struct {
	UserName     string `form:"username" check:"NotNull"`
	UserPassword string `form:"password"`
	Age          int    `form:"age" check:"Size(2);Range(12,45);NotNull"`
	Email        string `form:"email" check:"Email"`
}
