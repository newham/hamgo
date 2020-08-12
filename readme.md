# hamgo
>**hamgo** —— A tiny MVC web framework based on golang !   
You will find that build a website is **so easy** by using hamgo !  
Try it right now!  

**github** : [https://github.com/newham/hamgo](https://github.com/newham/hamgo)
## Getting Started
```go
go get github.com/newham/hamgo
```


## A simplest example
main.go
```go
package main

import (
    "github.com/newham/hamgo"
)

func main() {
    server := hamgo.New().Server()
    server.Get("/hello", Hello)
    server.RunAt("8080")
}

func Hello(ctx *hamgo.WebContext) {
    ctx.WriteString("Hello World!")
    ctx.Text(200)
}

```
then run it
```go
go run main.go
```
You will see at [http://localhost:8080/hello](http://localhost:8080/hello)
```
Hello world!
```
## Filter
main.go
```go
package main

import (
    "github.com/newham/hamgo"
)

func main() {
    server := hamgo.New().Server()
    server.Filter(LoginFilter).AddAnnoURL("/login")//set filter func LoginFilter(),and add anno url "/login";
    //Filter must be set before set handler ; 
    //Filter can be set only one in this version
    server.Get("/hello",Hello)
    server.RunAt("8080")
}

//return: true is pass filter , false is not pass
func LoginFilter(ctx *hamgo.WebContext) bool {
    if ctx.GetSession().Get(USER_SESSION) != nil {
		return true
	}
	ctx.PutData("code", 401)
	ctx.PutData("msg", "Login please!")
	ctx.JSON(401)
	return false
}

func Hello(ctx *hamgo.WebContext) {
    ctx.WriteString("Hello World!")
    ctx.Text(200)
}
```
then run it
```go
go run main.go
```
You will see at [http://localhost:8080/hello](http://localhost:8080/hello)
```
{
    "msg":"Login please!"
}
```
not
```
Hello World!
```

## Controller AOP
main.go
```go
package main

import (
    "github.com/newham/hamgo"
)

func main() {
    server := hamgo.New(nil)
    server.Get("/hello", BeforeHello, Hello)
    server.RunAt("8080")
}

func Hello(ctx *hamgo.WebContext) {
    ctx.WriteString("Hello World!")
    ctx.Text(200)
}

```
then run it
```go
go run main.go
```
You will see at [http://localhost:8080/hello](http://localhost:8080/hello)
```
Hello World!
```

## HTML & Template As Response  
main.go
```go
package main

import (
    "github.com/newham/hamgo"
)

func main() {
    server := hamgo.New(nil)
    server.Get("/page", Page)
    server.RunAt("8080")
}

func Page(ctx *hamgo.WebContext) {
    ctx.PutData("title", "This is title")
    ctx.HTML("index.html", "index_title.tmpl") // ctx.HTML( [html],[template1],[template2]...) 
}

```
index_title.tmpl
```html
{{ define "IndexTitle"}}
    <h1>{{.title}}</h1>
{{end}}
```
index.html
```html
<html>
    <head></head>
    <body>
        {{template "IndexTitle" .}} <!-- Do not forget '.' -->
        This is body
    </body>
</html>
```
then run it
```go
go run main.go
```
You will see at [http://localhost:8080/page](http://localhost:8080/page)
```
This is title
This is body
```
more doc about template of golang at: [golang template](https://newham.github.io/hamgo-doc/)

## Json As Response  
```go
//return json by putData
func Json(ctx *hamgo.WebContext) {
	ctx.PutData("say", "hello world")
	ctx.JSON(200)
}

//return json by yourself data ( example:map[string]interface{} )
func JsonFromData(ctx *hamgo.WebContext) {
	data := map[string]interface{}{"say": "hello world"}
	ctx.JSONFrom(200, data)
}

```
return json result
```json
{
    "say":"hello world"
}
```
## Static File
code
```go
server := hamgo.New().Server()
server.Static("public")
```
html
```html
<link rel="icon" href="/public/img/logo.ico" type="image/x-icon" />
<script src="public/js/jquery.min.js"></script>
```

## Restful API
set restful API controller , [start with '*=*']
```go
server.Get("/index/=model/=id", controller.Index)
```
controller
```go
func Index(ctx *hamgo.WebContext) {
	model := ctx.PathParam("model")
    id := ctx.PathParam("id")
}
```

## Config
### [1] init
init config file at create server
```go
server := hamgo.NewUseConf("./app.conf").Server()
```
set config file after create server
```go
server := hamgo.New().UseConfig("./app.conf").Server()
```
app.conf ( config file )
```conf
index = "hello"

port = 8087

# second
session_max_time = 1800

```
### [2] use
```go
port := hamgo.Conf.String("port")
```

## Logger
### [1] init
```go
server := hamgo.UseConfig("./log/app.log").Server()
```

### [2] config
app.conf ( config file )
```go
port = 8087

# second
session_max_time = 1800

# logger
log_console = true
log_file = "log/app.log"
# KB
log_file_max_size = 50
# KB
log_buf_size = 10
# s
log_buf_time = 10
# log format, example : [Error] 2018-01-01 14:35:16 /xxxx/test.go:75 -> 401,Unauthorized
log_format = "[%Title] %Time %File : %Text"
```

### [3] use
```go
hamgo.Log.Debug("done old UserName:%s", user.UserName)
hamgo.Log.Warn("UserPassword:%s", user.UserPassword)
hamgo.Log.Info("Age:%d", user.Age)
hamgo.Log.Error("Email:%s", user.Email)
```
you will see output in [./log/app.log] and [console]
```bash
[Debug] [2017-06-09 17:06:49] [test.go:55] done old UserName:test_user
[Warn] [2017-06-09 17:06:49] [test.go:56] UserPassword:123
[Info] [2017-06-09 17:06:49] [test.go:57] Age:23
[Error] [2017-06-09 17:06:49] [test.go:58] Email:test@test.com
```

## Session
### [1] init
```go
server := hamgo.UseSession(3600).Server() //session timeout is 3600 seconds

```
### [2] use
```go
type User struct{
    UserName string
    Password string
}

var user User
session :=ctx.GetSession() //get session
session.Set("user",user) //set session key-value
user := session.Get("user").(User) //get session value by key
session.Delete("user") //delete session value by key
sessionId :=session.SessionID() //get session id
```

## Bind Request Form or Json & do validate
```go
type MyForm struct {
    UserName     string `form:"username" check:"NotNull"`
    UserPassword string `form:"password" check:"MinSize(8);NotNull"`
    Age          int    `form:"age" check:"Range(12,45);NotNull"`
    Email        string `form:"email" check:"Email"`
}

func Bind(ctx *hamgo.WebContext) {
    form := MyForm{}
    errs := ctx.BindForm(&form) // Do not forget '&'

    //print validate errors
    if len(errs) > 0 {
        for k, err := range errs {
            hamgo.Log.Error("%s:%s", k, err.Error())
        }
    }

    //after bind , you can use your binded form or json data
    println(form.UserName)
}
```

## What hamgo have ?

|Features      |Support  |
|--------------|:-------:|
|restful API   |√        |
|filter        |√        |
|handler AOP   |√        |
|config        |√        |
|session       |√        |
|logger        |√        |
|form validate |√        |
|response json |√        |
|json,form bind|√        |
|html,template |√        |
More features are coming soon ... 

## More Documentation
[https://newham.github.io/hamgo-doc/](https://newham.github.io/hamgo-doc/)

## License
This project is under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html)
```
Copyright newham.cn@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

## Feedback
If you have any question or suggestion , and if you meet some **bug** ,send e-mail to me at [**newham.cn@gmail.com**]