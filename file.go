package hamgo

import (
	"os"
	"path/filepath"
	"strings"
)

func writeString(filename string, content string) bool {
	f := openFile(filename)
	if f == nil {
		return false
	}
	defer f.Close()
	_, err := f.WriteString(content)
	if err != nil {
		println("append file failed!", err.Error())
		return false
	}

	return true
}

func writeBytes(filename string, content []byte) bool {
	f := openFile(filename)
	if f == nil {
		return false
	}
	defer f.Close()
	_, err := f.Write(content)
	if err != nil {
		println("append file failed!", err.Error())
		return false
	}

	return true
}

func isFileExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func deleteFile(filename string) bool {
	err := os.RemoveAll(filename)
	if err != nil {
		return false
	}
	return true
}

func openFile(filename string) *os.File {
	var f *os.File
	var err error
	if !isFileExist(filename) {
		err = os.MkdirAll(filepath.Dir(filename), 0755)
		if err != nil {
			println("mk dir failed ", filename, " failed,", err)
			return nil
		}
	}
	f, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		println("open file failed ", filename, " failed,", err)
		return nil
	}
	return f
}

func renameFile(filename, newname string) bool {
	err := os.Rename(filename, newname)
	if err == nil {
		return true
	}
	return false
}

func currentPath(filename string) string {
	index := strings.LastIndex(filename, "/")
	return filename[:index+1]

}
