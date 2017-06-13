# hamgo
### A tiny MVC web framework based on golang!
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
Hello world!
```

### HTML & Template
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
```
{{ define "IndexTitle"}}
    <h1>{{.title}}</h1>
{{end}}
```
index.html
```
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

### Bind Request Form & Json
...

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
|json,form bind|√        |
|html,template |√        |
More features are coming soon ... 

## More Documentation
[https://newham.github.io/hamgo-doc/](https://newham.github.io/hamgo-doc/)