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
	JSON(code int) error
	JSONFrom(code int, data interface{}) error
	JSONString(code int, data string) error
	DataHTML(data interface{}, filenames ...string)
	Redirect(code int, path string)
	Code(statusCode int)
	PathParam(key string) string
	FormFile(fileName string) (multipart.File, *multipart.FileHeader, error)
	GetSession() Session
	BindForm(obj interface{}) map[string]error
	BindJSON(obj interface{}) error
	HTML(filenames ...string)
	PutData(key string, data interface{})
}

//WebContext :
type WebContext struct {
	Context
	W          http.ResponseWriter
	R          *http.Request
	RespBuf    *bytes.Buffer
	StatusCode int
	PathParams map[string]string
	Data       map[string]interface{}
}

//NewWebContext :
func newWebContext(rw http.ResponseWriter, r *http.Request, path string) *WebContext {
	return &WebContext{W: rw, R: r,
		RespBuf:    new(bytes.Buffer),
		StatusCode: http.StatusOK,
		PathParams: newPath(path).PathParam(r.URL.Path),
		Data:       make(map[string]interface{})}
}

//FormValue :
func (ctx *WebContext) FormValue(key string) string {
	return ctx.R.FormValue(key)
}

//WriteBytes :
func (ctx *WebContext) WriteBytes(b []byte) {
	ctx.RespBuf.Write(b)
}

//WriteString :
func (ctx *WebContext) WriteString(str string) {
	ctx.RespBuf.WriteString(str)
}

//Text :
func (ctx *WebContext) Text(code int) {
	ctx.Code(code)
	ctx.W.WriteHeader(ctx.StatusCode)
	ctx.W.Write(ctx.RespBuf.Bytes())
}

//JSON :
func (ctx *WebContext) JSON(code int) error {
	ctx.Code(code)
	ctx.W.WriteHeader(ctx.StatusCode)
	b, err := json.Marshal(ctx.Data)
	if err != nil {
		return err
	}
	_, err = ctx.W.Write(b)
	return err
}

//JSONFrom :
func (ctx *WebContext) JSONFrom(code int, data interface{}) error {
	ctx.Code(code)
	ctx.W.WriteHeader(ctx.StatusCode)
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = ctx.W.Write(b)
	return err
}

//JSONString :
func (ctx *WebContext) JSONString(code int, data string) error {
	ctx.Code(code)
	ctx.W.WriteHeader(ctx.StatusCode)
	_, err := ctx.W.Write([]byte(data))
	return err
}

//DataHTML :
func (ctx *WebContext) DataHTML(data interface{}, filenames ...string) {
	t, err := template.ParseFiles(filenames...)
	if err != nil {
		ctx.WriteString("prase template failed! check file path")
		ctx.Text(500)
		return
	}
	t.Execute(ctx.W, data)
}

//Redirect :
func (ctx *WebContext) Redirect(code int, path string) {
	http.Redirect(ctx.W, ctx.R, path, code)
}

//Code :
func (ctx *WebContext) Code(statusCode int) {
	ctx.StatusCode = statusCode
}

//PathParam :
func (ctx *WebContext) PathParam(key string) string {
	return ctx.PathParams[key]
}

//FormFile :
func (ctx *WebContext) FormFile(fileName string) (multipart.File, *multipart.FileHeader, error) {
	return ctx.R.FormFile(fileName)
}

//GetSession :
func (ctx *WebContext) GetSession() Session {
	return sessions.SessionStart(ctx.W, ctx.R)
}

//BindForm : use reflect to bind form-values to object
func (ctx *WebContext) BindForm(obj interface{}) map[string]error {
	errs := make(map[string]error)
	rt := reflect.TypeOf(obj).Elem()
	rv := reflect.ValueOf(obj).Elem()
	for i := 0; i < rt.NumField(); i++ {
		rf := rt.Field(i)
		formName := rf.Tag.Get("form")
		formValue := ctx.R.FormValue(formName)
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
func (ctx *WebContext) BindJSON(obj interface{}) error {
	j, err := ioutil.ReadAll(ctx.R.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, obj)
}

//HTML :
func (ctx *WebContext) HTML(filenames ...string) {
	ctx.DataHTML(ctx.Data, filenames...)
}

//PutData :
func (ctx *WebContext) PutData(key string, data interface{}) {
	ctx.Data[key] = data
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
			if m, _ := regexp.MatchString(`^(\(\d{3,4}-)|\d{3.4}-)?\d{7,8}$`, formValue); !m {
				return newError(checkTagErrorTel)
			}
		case checkTagPhone:
			m1, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, formValue)
			m2, _ := regexp.MatchString(`^(\(\d{3,4}-)|\d{3.4}-)?\d{7,8}$`, formValue)
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
