package hamgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

const (
	checkSplit = ";"
	//tag
	checkTagEmail   = "Email"
	checkTagTel     = "Tel"
	checkTagPhone   = "Phone"
	checkTagURL     = "Url"
	checkTagNum     = "Num"
	checkTagMobile  = "Mobile"
	checkTagIpv4    = "Ipv4"
	checkTagSize    = "Size"
	checkTagMaxSize = "MaxSize"
	checkTagMinSize = "MinSize"
	checkTagNotNull = "NotNull"
	checkTagMin     = "Min"
	checkTagMax     = "Max"
	checkTagRange   = "Range"
	//error
	checkTagErrorEmail   = "wrong email"
	checkTagErrorTel     = "wrong tel"
	checkTagErrorPhone   = "wrong phone"
	checkTagErrorURL     = "wrong url"
	checkTagErrorNum     = "wrong number"
	checkTagErrorMobile  = "wrong mobile"
	checkTagErrorIpv4    = "wrong ipv4"
	checkTagErrorSize    = "wrong size"
	checkTagErrorMaxSize = "wrong maxsize"
	checkTagErrorMinSize = "wrong minsize"
	checkTagErrorNotNull = "null value"
	checkTagErrorMin     = "wrong min"
	checkTagErrorMax     = "wrong max"
	checkTagErrorRange   = "wrong range"
)

type Context interface {
	FormValue(key string) string
	WriteBytes(b []byte)
	WriteString(str string)
	Text(code int)
	JSON(code int, b []byte) error
	JSONFrom(code int, data interface{}) error
	JSONString(code int, data string) error
	DataHTML(data interface{}, filenames ...string)
	Redirect(path string)
	Code(statusCode int)
	PathParam(key string) string
	FormFile(fileName string) (multipart.File, *multipart.FileHeader, error)
	GetSession() Session
	DeleteSession()
	BindForm(obj interface{}) map[string]error
	BindJSON(obj interface{}) error
	HTML(filenames ...string)
	File(filename string)
	PutData(key string, data interface{})
	R() *http.Request
	W() http.ResponseWriter
	Method() string
	Body() ([]byte, error)
}

//webContext :
type webContext struct {
	Context
	w          http.ResponseWriter
	r          *http.Request
	respBuf    *bytes.Buffer
	statusCode int
	pathParams map[string]string
	data       map[string]interface{}
}

//NewWebContext :
func newWebContext(rw http.ResponseWriter, r *http.Request, path string) Context {
	return &webContext{w: rw, r: r,
		respBuf:    new(bytes.Buffer),
		statusCode: http.StatusOK,
		pathParams: newPath(path).PathParam(r.URL.Path),
		data:       make(map[string]interface{})}
}

//FormValue :
func (ctx *webContext) FormValue(key string) string {
	return ctx.r.FormValue(key)
}

//WriteBytes :
func (ctx *webContext) WriteBytes(b []byte) {
	ctx.respBuf.Write(b)
}

//WriteString :
func (ctx *webContext) WriteString(str string) {
	ctx.respBuf.WriteString(str)
}

//Text :
func (ctx *webContext) Text(code int) {
	ctx.Code(code)
	ctx.w.WriteHeader(ctx.statusCode)
	ctx.w.Write(ctx.respBuf.Bytes())
}

//JSON :
func (ctx *webContext) JSON(code int, b []byte) error {
	ctx.Code(code)
	ctx.w.WriteHeader(ctx.statusCode)
	_, err := ctx.w.Write(b)
	return err
}

//JSONFrom :
func (ctx *webContext) JSONFrom(code int, data interface{}) error {
	ctx.Code(code)
	ctx.w.WriteHeader(ctx.statusCode)
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = ctx.w.Write(b)
	return err
}

//JSONString :
func (ctx *webContext) JSONString(code int, data string) error {
	ctx.Code(code)
	ctx.w.WriteHeader(ctx.statusCode)
	_, err := ctx.w.Write([]byte(data))
	return err
}

//DataHTML :
func (ctx *webContext) DataHTML(data interface{}, filenames ...string) {
	t, err := template.ParseFiles(filenames...)
	if err != nil {
		ctx.WriteString("prase template failed! check file path")
		ctx.Text(500)
		return
	}
	t.Execute(ctx.w, data)
}

//Redirect :
func (ctx *webContext) Redirect(path string) {
	ctx.w.Header().Set("Location", path)
	ctx.w.WriteHeader(301)
}

//Code :
func (ctx *webContext) Code(statusCode int) {
	ctx.statusCode = statusCode
}

//PathParam :
func (ctx *webContext) PathParam(key string) string {
	return ctx.pathParams[key]
}

//FormFile :
func (ctx *webContext) FormFile(fileName string) (multipart.File, *multipart.FileHeader, error) {
	return ctx.r.FormFile(fileName)
}

//GetSession :
func (ctx *webContext) GetSession() Session {
	return sessions.SessionStart(ctx.w, ctx.r)
}

//DeleteSession :
func (ctx *webContext) DeleteSession() {
	sessions.SessionDestroy(ctx.w, ctx.r)
}

//BindForm : use reflect to bind form-values to object
func (ctx *webContext) BindForm(obj interface{}) map[string]error {
	errs := make(map[string]error)
	rt := reflect.TypeOf(obj).Elem()
	rv := reflect.ValueOf(obj).Elem()
	for i := 0; i < rt.NumField(); i++ {
		rf := rt.Field(i)
		formName := rf.Tag.Get("form")
		formValue := ctx.r.FormValue(formName)
		//1.check value
		if err := checkValueByTag(formName, formValue, rf.Tag.Get("check")); err != nil {
			errs[formName] = err
		}
		//2.set value
		switch rf.Type.Kind() {
		case reflect.String:
			rv.Field(i).SetString(formValue)
		case reflect.Int:
			rv.Field(i).SetInt(stringToInt64(formValue, 0))
		case reflect.Int64:
			rv.Field(i).SetInt(stringToInt64(formValue, 0))
		case reflect.Float32:
			rv.Field(i).SetFloat(stringToFloat64(formValue, 0))
		case reflect.Float64:
			rv.Field(i).SetFloat(stringToFloat64(formValue, 0))
		}
	}
	return errs
}

//BindJSON : prase JSON to Struct
func (ctx *webContext) BindJSON(obj interface{}) error {
	j, err := ioutil.ReadAll(ctx.r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, obj)
}

//HTML :
func (ctx *webContext) HTML(filenames ...string) {
	ctx.DataHTML(ctx.data, filenames...)
}

//File :
func (ctx *webContext) File(filename string) {
	http.ServeFile(ctx.w, ctx.r, filename)
}

//PutData :
func (ctx *webContext) PutData(key string, data interface{}) {
	ctx.data[key] = data
}

//R() :
func (ctx *webContext) R() *http.Request {
	return ctx.r
}

//W() :
func (ctx *webContext) W() http.ResponseWriter {
	return ctx.w
}

//Method() :
func (ctx *webContext) Method() string {
	return ctx.r.Method
}

//Body() :
func (ctx *webContext) Body() ([]byte, error) {
	return ioutil.ReadAll(ctx.r.Body)
}

func checkValueByTag(formName, formValue, check string) error {
	tags := strings.Split(check, checkSplit)
	for _, tag := range tags {
		if tag == "" {
			continue
		}
		tagStart := strings.Index(tag, "(")
		tagEnd := strings.Index(tag, ")")
		tagVal := ""
		tagName := tag
		if tagStart != -1 {
			tagName = tag[:tagStart]
			tagVal = tag[tagStart+1 : tagEnd]
		}
		switch tagName {
		case checkTagEmail:
			if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, formValue); !m {
				return newError(checkTagErrorEmail)
			}
		case checkTagMobile:
			if m, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, formValue); !m {
				return newError(checkTagErrorMobile)
			}
		case checkTagNotNull:
			if formValue == "" {
				return newError(checkTagErrorNotNull)
			}
		case checkTagTel:
			if m, _ := regexp.MatchString(`^(\(\d{3,4}-)|\d{3.4}-)`+`?\d{7,8}$`, formValue); !m {
				return newError(checkTagErrorTel)
			}
		case checkTagPhone:
			m1, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, formValue)
			m2, _ := regexp.MatchString(`^(\(\d{3,4}-)|\d{3.4}-)`+`?\d{7,8}$`, formValue)
			if !m1 && !m2 {
				return newError(checkTagErrorPhone)
			}
		case checkTagURL:
			if m, _ := regexp.MatchString(`[a-zA-z]+://[^\s]*`, formValue); !m {
				return newError(checkTagErrorURL)
			}
		case checkTagNum:
			if m, _ := regexp.MatchString(`^[0-9]*$`, formValue); !m {
				return newError(checkTagErrorNum)
			}
		case checkTagIpv4:
			if m, _ := regexp.MatchString(`((?:(?:25[0-5]|2[0-4]\\d|[01]?\\d?\\d)\\.){3}(?:25[0-5]|2[0-4]\\d|[01]?\\d?\\d))`, formValue); !m {
				return newError(checkTagErrorIpv4)
			}
		case checkTagSize:
			if m, _ := regexp.MatchString(`^.{`+tagVal+`}$`, formValue); !m {
				return newError(checkTagErrorSize)
			}
		case checkTagMaxSize:
			if m, _ := regexp.MatchString(`^.{,`+tagVal+`}$`, formValue); !m {
				return newError(checkTagErrorMaxSize)
			}
		case checkTagMinSize:
			if m, _ := regexp.MatchString(`^.{`+tagVal+`,}$`, formValue); !m {
				return newError(checkTagErrorMinSize)
			}
		case checkTagMin:
			if stringToInt(formValue, 0) < stringToInt(tagVal, 0) {
				return newError(checkTagErrorMin)
			}
		case checkTagMax:
			if stringToInt(formValue, 0) > stringToInt(tagVal, 0) {
				return newError(checkTagErrorMin)
			}
		case checkTagRange:
			vals := strings.Split(tagVal, ",")
			min := stringToInt(vals[0], 0)
			max := stringToInt(vals[1], 1)
			val := stringToInt(formValue, 0)
			if min > val || max < val {
				return newError(checkTagErrorRange)
			}
		}
	}

	return nil
}

func newError(errorInfo string) error {
	return fmt.Errorf(errorInfo)
}
