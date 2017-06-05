package hamgo

import (
	"bytes"
	"encoding/json"
	"html/template"
	"mime/multipart"
	"net/http"
	"reflect"
)

//WebContext :
type WebContext struct {
	W          http.ResponseWriter
	R          *http.Request
	RespBuf    *bytes.Buffer
	StatusCode int
	PathParams map[string]string
}

//NewWebContext :
func newWebContext(rw http.ResponseWriter, r *http.Request, path string) WebContext {
	return WebContext{W: rw, R: r, RespBuf: new(bytes.Buffer), StatusCode: http.StatusOK, PathParams: newPath(path).PathParam(r.URL.Path)}
}

//FormValue :
func (ctx WebContext) FormValue(key string) string {
	return ctx.R.FormValue(key)
}

//WriteBytes :
func (ctx WebContext) WriteBytes(b []byte) {
	ctx.RespBuf.Write(b)
}

//WriteString :
func (ctx WebContext) WriteString(str string) {
	ctx.RespBuf.WriteString(str)
}

//Text :
func (ctx WebContext) Text(code int) {
	ctx.Code(code)
	ctx.W.WriteHeader(ctx.StatusCode)
	ctx.W.Write(ctx.RespBuf.Bytes())
}

//Json :
func (ctx WebContext) Json(code int) error {
	ctx.Code(code)
	ctx.W.WriteHeader(ctx.StatusCode)
	b, err := json.Marshal(ctx.RespBuf.Bytes())
	if err != nil {
		return err
	}
	ctx.W.Write(b)
	return nil
}

//JsonFrom :
func (ctx WebContext) JsonFrom(code int, data interface{}) error {
	ctx.Code(code)
	ctx.W.WriteHeader(ctx.StatusCode)
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	ctx.W.Write(b)
	return nil
}

//Html :
func (ctx WebContext) Html(data interface{}, filenames ...string) {
	t, err := template.ParseFiles(filenames...)
	if err != nil {
		ctx.WriteString("prase template failed! check file path")
		ctx.Text(500)
		return
	}
	t.Execute(ctx.W, data)
}

//Redirect :
func (ctx WebContext) Redirect(code int, path string) {
	http.Redirect(ctx.W, ctx.R, path, code)
}

//Code :
func (ctx WebContext) Code(statusCode int) {
	ctx.StatusCode = statusCode
}

//PathParam :
func (ctx WebContext) PathParam(key string) string {
	return ctx.PathParams[key]
}

//FormFile :
func (ctx WebContext) FormFile(fileName string) (multipart.File, *multipart.FileHeader, error) {
	return ctx.R.FormFile(fileName)
}

//GetSession :
func (ctx WebContext) GetSession() Session {
	return Sessions.SessionStart(ctx.W, ctx.R)
}

//BindForm
func (ctx WebContext) BindForm(obj interface{}) interface{} {
	s := reflect.TypeOf(obj)
	// sv := reflect.ValueOf(obj)
	svr := reflect.New(reflect.ValueOf(obj).Type())
	for i := 0; i < s.NumField(); i++ {
		sf := s.Field(i)
		switch sf.Type.String() {
		case "string":
			formName := sf.Tag.Get("form")
			println("formName:" + formName)
			formValue := ctx.R.FormValue(formName)
			println("formValue:" + formValue)

			svr.Elem().Field(i).SetString(formValue)

		}
	}
	return svr.Interface()
}
