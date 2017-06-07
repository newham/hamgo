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
	CheckSplit = ";"
	//tag
	CheckTagEmail   = "Email"
	CheckTagTel     = "Tel"
	CheckTagPhone   = "Phone"
	CheckTagUrl     = "Url"
	CheckTagNum     = "Num"
	CheckTagMobile  = "Mobile"
	CheckTagIpv4    = "Ipv4"
	CheckTagSize    = "Size"
	CheckTagMaxSize = "MaxSize"
	CheckTagMinSize = "MinSize"
	CheckTagNotNull = "NotNull"
	CheckTagMin     = "Min"
	CheckTagMax     = "Max"
	CheckTagRange   = "Range"
	//error
	CheckTagErrorEmail   = "wrong email"
	CheckTagErrorTel     = "wrong tel"
	CheckTagErrorPhone   = "wrong phone"
	CheckTagErrorUrl     = "wrong url"
	CheckTagErrorNum     = "wrong number"
	CheckTagErrorMobile  = "wrong mobile"
	CheckTagErrorIpv4    = "wrong ipv4"
	CheckTagErrorSize    = "wrong size"
	CheckTagErrorMaxSize = "wrong maxsize"
	CheckTagErrorMinSize = "wrong minsize"
	CheckTagErrorNotNull = "null value"
	CheckTagErrorMin     = "wrong min"
	CheckTagErrorMax     = "wrong max"
	CheckTagErrorRange   = "wrong range"
)

//WebContext :
type WebContext struct {
	W          http.ResponseWriter
	R          *http.Request
	RespBuf    *bytes.Buffer
	StatusCode int
	PathParams map[string]string
	Data       map[string]interface{}
}

//NewWebContext :
func newWebContext(rw http.ResponseWriter, r *http.Request, path string) WebContext {
	return WebContext{W: rw, R: r,
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

//Json :
func (ctx *WebContext) Json(code int) error {
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
func (ctx *WebContext) JsonFrom(code int, data interface{}) error {
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
func (ctx *WebContext) Html(data interface{}, filenames ...string) {
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
	return Sessions.SessionStart(ctx.W, ctx.R)
}

//BindForm : use reflect to bind form value to object
func (ctx *WebContext) BindForm(obj interface{}) error {
	rt := reflect.TypeOf(obj).Elem()
	rv := reflect.ValueOf(obj).Elem()
	for i := 0; i < rt.NumField(); i++ {
		rf := rt.Field(i)
		formName := rf.Tag.Get("form")
		formValue := ctx.R.FormValue(formName)
		//1.check value
		if err := checkValueByTag(formName, formValue, rf.Tag.Get("check")); err != nil {
			return err
		}
		//2.set value
		switch rf.Type.Kind() {
		case reflect.String:
			rv.Field(i).SetString(formValue)
		case reflect.Int:
			rv.Field(i).SetInt(StringToInt64(formValue, 0))
		case reflect.Int64:
			rv.Field(i).SetInt(StringToInt64(formValue, 0))
		case reflect.Float32:
			rv.Field(i).SetFloat(StringToFloat64(formValue, 0))
		case reflect.Float64:
			rv.Field(i).SetFloat(StringToFloat64(formValue, 0))
		}
	}
	return nil
}

func (ctx *WebContext) BindJson(obj interface{}) error {
	j, err := ioutil.ReadAll(ctx.R.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, obj)
}

//Form :
func (ctx *WebContext) DataHtml(filenames ...string) {
	ctx.Html(ctx.Data, filenames...)
}

func checkValueByTag(formName, formValue, check string) error {
	tags := strings.Split(check, CheckSplit)
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
		case CheckTagEmail:
			if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, formValue); !m {
				return newError(formName, CheckTagErrorEmail)
			}
		case CheckTagMobile:
			if m, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, formValue); !m {
				return newError(formName, CheckTagErrorMobile)
			}
		case CheckTagNotNull:
			if formValue == "" {
				return newError(formName, CheckTagErrorNotNull)
			}
		case CheckTagTel:
			if m, _ := regexp.MatchString(`^(\(\d{3,4}-)|\d{3.4}-)?\d{7,8}$`, formValue); !m {
				return newError(formName, CheckTagErrorTel)
			}
		case CheckTagPhone:
			if m, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})|(\(\d{3,4}-)|\d{3.4}-)?\d{7,8}$`, formValue); !m {
				return newError(formName, CheckTagErrorPhone)
			}
		case CheckTagUrl:
			if m, _ := regexp.MatchString(`[a-zA-z]+://[^\s]*`, formValue); !m {
				return newError(formName, CheckTagErrorUrl)
			}
		case CheckTagNum:
			if m, _ := regexp.MatchString(`^[0-9]*$`, formValue); !m {
				return newError(formName, CheckTagErrorNum)
			}
		case CheckTagIpv4:
			if m, _ := regexp.MatchString(`((?:(?:25[0-5]|2[0-4]\\d|[01]?\\d?\\d)\\.){3}(?:25[0-5]|2[0-4]\\d|[01]?\\d?\\d))`, formValue); !m {
				return newError(formName, CheckTagErrorIpv4)
			}
		case CheckTagSize:
			if m, _ := regexp.MatchString(`^.{`+tagVal+`}$`, formValue); !m {
				return newError(formName, CheckTagErrorSize)
			}
		case CheckTagMaxSize:
			if m, _ := regexp.MatchString(`^.{,`+tagVal+`}$`, formValue); !m {
				return newError(formName, CheckTagErrorMaxSize)
			}
		case CheckTagMinSize:
			if m, _ := regexp.MatchString(`^.{`+tagVal+`,}$`, formValue); !m {
				return newError(formName, CheckTagErrorMinSize)
			}
		case CheckTagMin:
			if StringToInt(formValue, 0) < StringToInt(tagVal, 0) {
				return newError(formName, CheckTagErrorMin)
			}
		case CheckTagMax:
			if StringToInt(formValue, 0) > StringToInt(tagVal, 0) {
				return newError(formName, CheckTagErrorMin)
			}
		case CheckTagRange:
			vals := strings.Split(tagVal, ",")
			min := StringToInt(vals[0], 0)
			max := StringToInt(vals[1], 1)
			val := StringToInt(formValue, 0)
			if min > val || max < val {
				return newError(formName, CheckTagErrorRange)
			}
		}
	}

	return nil
}
func newError(formName, errorInfo string) error {
	return fmt.Errorf("%s:%s", formName, errorInfo)
}
