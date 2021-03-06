package hamgo

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"strconv"
	"strings"
)

func int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func intToString(i int) string {
	return strconv.Itoa(i)
}

func stringToInt(str string, must int) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		return must
	}
	return value
}

func stringToInt64(str string, must int64) int64 {
	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return must
	}
	return value
}

func stringToFloat32(str string, must float32) float32 {
	value, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return must
	}
	return float32(value)
}

func stringToFloat64(str string, must float64) float64 {
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return must
	}
	return value
}

func getFileName(path string) string {
	return path[strings.LastIndex(path, "/")+1:]
}

func uuid(len int) string {
	b := make([]byte, len)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func JSONToString(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(b)
}
