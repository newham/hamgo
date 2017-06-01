package hamgo

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Context struct {
	W          http.ResponseWriter
	R          *http.Request
	RespBuf    *bytes.Buffer
	StatusCode int
}

func NewContext(rw http.ResponseWriter, r *http.Request) *Context {
	return &Context{rw, r, new(bytes.Buffer), http.StatusOK}
}

func (ctx *Context) FormValue(key string) string {
	return ctx.R.FormValue(key)
}

func (ctx *Context) WriteBytes(b []byte) {
	ctx.RespBuf.Write(b)
}

func (ctx *Context) WriteString(str string) {
	ctx.RespBuf.WriteString(str)
}

func (ctx *Context) Text(code int) {
	ctx.Code(code)
	ctx.W.WriteHeader(ctx.StatusCode)
	ctx.W.Write(ctx.RespBuf.Bytes())
}

func (ctx *Context) Json(code int) error {
	ctx.Code(code)
	ctx.W.WriteHeader(ctx.StatusCode)
	b, err := json.Marshal(ctx.RespBuf.Bytes())
	if err != nil {
		return err
	}
	ctx.W.Write(b)
	return nil
}

func (ctx *Context) JsonFrom(code int, data interface{}) error {
	ctx.Code(code)
	ctx.W.WriteHeader(ctx.StatusCode)
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	ctx.W.Write(b)
	return nil
}

func (ctx *Context) Html(view string) {

}
func (ctx *Context) Redirect(code int, path string) {
	http.Redirect(ctx.W, ctx.R, path, code)
}
func (ctx *Context) Code(statusCode int) {
	ctx.StatusCode = statusCode
}
