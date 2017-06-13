# hamgo
**hamgo** —— A tiny MVC web framework based on golang!   
You will find that build a website is **so easy** by using hamgo!  
Try it right now!
## Getting Started
```go
go get github.com/newham/hamgo
```

### A simplest example
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
You will see at [http://localhost:8080/hello](http://localhost:8080/hello)
```
Hello world!
```

### Handler Before & After
main.go
```go
package main

import (
    "github.com/newham/hamgo"
)

func main() {
    server := hamgo.New().Server()
    server.GetBefore("/hello", BeforeHello, Hello)
    server.RunAt("8080")
}

func Hello(ctx *hamgo.WebContext) {
    ctx.WriteString("Hello World!")
    ctx.Text(200)
}

func BeforeHello(ctx *hamgo.WebContext) {
    ctx.WriteString("Do Before Hello!\n")
}

```
You will see at [http://localhost:8080/hello](http://localhost:8080/hello)
```
Do Before Hello!
Hello World!
```

### Response HTML & Template
main.go
```go
package main

import (
    "github.com/newham/hamgo"
)

func main() {
    server := hamgo.New().Server()
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
You will see at [http://localhost:8080/page](http://localhost:8080/page)
```
This is title
This is body
```
more doc about template of golang at: [golang template](https://newham.github.io/hamgo-doc/)

### Response Json
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

### Bind Request Form or Json & do validate
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

## What we have ?

|Features      |Support  |
|--------------|:-------:|
|restful API   |√        |
|method filter |√        |
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