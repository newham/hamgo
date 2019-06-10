package hamgo

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

const (
	defaultConfFile = "conf/app.conf"
)

type configInterface interface {
	//Set(key, val string) error   // support section::key type in given key when using ini type.
	String(key string) string    // support section::key type in key string when using ini and json type; Int,Int64,Bool,Float,DIY are same.
	Strings(key string) []string //get string slice
	Int(key string) (int, error)
	Int64(key string) (int64, error)
	Bool(key string) (bool, error)
	Float(key string) (float64, error)
	DefaultString(key string, defaultval string) string      // support section::key type in key string when using ini and json type; Int,Int64,Bool,Float,DIY are same.
	DefaultStrings(key string, defaultval []string) []string //get string slice
	DefaultInt(key string, defaultval int) int
	DefaultInt64(key string, defaultval int64) int64
	DefaultBool(key string, defaultval bool) bool
	DefaultFloat(key string, defaultval float64) float64
}

type config struct {
	File string
	Keys map[string]string
}

//Conf : user can get config items
var Conf configInterface

func setConfig(configFile string) {
	if configFile == "" {
		configFile = defaultConfFile
	}
	if !isFileExist(configFile) {
		return
	}
	config := &config{configFile, make(map[string]string)}
	if err := config.Prase(); err != nil {
		panic(err)
	}
	Conf = config
}

func (c *config) Prase() error {
	isEnd := false

	f, err := os.Open(c.File)
	defer f.Close()
	if err != nil {
		return errors.New("Open file" + c.File + " failed")
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if line != "" {
				isEnd = true
			} else {
				break
			}
		}
		line = strings.TrimSpace(line)
		if isCommentOut(line) {
			continue
		}
		firstIndex := strings.Index(line, "=")

		if firstIndex < 1 {
			continue
		} else {
			c.Keys[strings.Trim(line[:firstIndex], "\" ")] = strings.Trim(line[firstIndex+1:], "\" ")
		}

		if isEnd {
			break
		}
	}
	return nil
}

func isCommentOut(line string) bool {
	if strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "*") || strings.HasPrefix(line, "[") {
		return true
	}
	return false

}

func (c *config) String(key string) string {
	return c.Keys[key]
}
func (c *config) Strings(key string) []string {
	if c.Keys[key] == "" {
		return make([]string, 0)
	}
	return strings.Split(c.Keys[key], " ")

}
func (c *config) Int(key string) (int, error) {
	return strconv.Atoi(c.Keys[key])
}
func (c *config) Int64(key string) (int64, error) {
	return strconv.ParseInt(c.Keys[key], 10, 64)
}
func (c *config) Bool(key string) (bool, error) {
	return strconv.ParseBool(c.Keys[key])
}
func (c *config) Float(key string) (float64, error) {
	return strconv.ParseFloat(c.Keys[key], 64)
}

func (c *config) DefaultString(key string, defaultval string) string {
	if c.String(key) == "" {
		return defaultval
	}
	return c.String(key)

}
func (c *config) DefaultStrings(key string, defaultval []string) []string {
	if len(c.Strings(key)) < 1 {
		return defaultval
	}
	return c.Strings(key)

}
func (c *config) DefaultInt(key string, defaultval int) int {
	value, err := c.Int(key)
	if err != nil {
		return defaultval
	}
	return value

}
func (c *config) DefaultInt64(key string, defaultval int64) int64 {
	value, err := c.Int64(key)
	if err != nil {
		return defaultval
	}
	return value

}
func (c *config) DefaultBool(key string, defaultval bool) bool {
	value, err := c.Bool(key)
	if err != nil {
		return defaultval
	}
	return value

}
func (c *config) DefaultFloat(key string, defaultval float64) float64 {
	value, err := c.Float(key)
	if err != nil {
		return defaultval
	}
	return value

}
