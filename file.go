package hamgo

import (
	"os"
	"path/filepath"
)

func WriteString(filename string, content string) bool {
	f := OpenFile(filename)
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

func WriteBytes(filename string, content []byte) bool {
	f := OpenFile(filename)
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

func IsFileExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func DeleteFile(filename string) bool {
	err := os.RemoveAll(filename)
	if err != nil {
		return false
	}
	return true
}

func OpenFile(filename string) *os.File {
	var f *os.File
	var err error
	if !IsFileExist(filename) {
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

func FileSize(filename string) int64 {
	f := OpenFile(filename)
	if f == nil {
		return -1
	}
	defer f.Close()

	fs, err := f.Stat()
	if err != nil {
		return -1
	}

	return fs.Size()
}

func RenameFile(filename, newname string) bool {
	err := os.Rename(filename, newname)
	if err == nil {
		return true
	}
	return false
}
