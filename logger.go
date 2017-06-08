package hamgo

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"
)

var Log Logger

type Logger interface {
	Error(format string, a ...interface{})
	Info(format string, a ...interface{})
	Debug(format string, a ...interface{})
	Warn(format string, a ...interface{})
	WriteBuf()
}

type FileLogger struct {
	FilePath    string
	fileMaxSize int64
	Format      string
	buf         *bytes.Buffer
	mutex       *sync.Mutex
	bufTime     time.Duration
	bufSize     int
	console     bool
}

const (
	DefaultWriteBufTime = 2000      //2000ms
	DefaultWriteBufSize = 1 * 1024  //1KB
	DefaultFileMaxSize  = 10 * 1024 //10KB
	LogTitleDebug       = "Debug"
	LogTitleInfo        = "Info"
	LogTitleError       = "Error"
	LogTitleWarn        = "Warn"
)

func newLogger(filePath string) {
	Log = &FileLogger{
		FilePath:    filePath,
		console:     true,
		buf:         new(bytes.Buffer),
		mutex:       new(sync.Mutex),
		fileMaxSize: DefaultFileMaxSize,
		bufTime:     DefaultWriteBufTime,
		bufSize:     DefaultWriteBufSize}
	//create a thread to write buf to log file
	go Log.WriteBuf()
}

func newLoggerFromConf() {
	Log = &FileLogger{
		console: true,
		buf:     new(bytes.Buffer),
		mutex:   new(sync.Mutex),
		bufTime: DefaultWriteBufTime,
		bufSize: DefaultWriteBufSize}
	//create a thread to write buf to log file
	go Log.WriteBuf()
}

func (log *FileLogger) Error(format string, a ...interface{}) {
	log.writeAndPrint("Error", format, a...)
}

func (log *FileLogger) Info(format string, a ...interface{}) {
	log.writeAndPrint("Info", format, a...)
}

func (log *FileLogger) Debug(format string, a ...interface{}) {
	log.writeAndPrint("Debug", format, a...)
}

func (log *FileLogger) Warn(format string, a ...interface{}) {
	log.writeAndPrint("Warn", format, a...)
}

func (log *FileLogger) WriteBuf() {
	//listen exit signal
	log.onExit()
	for {
		//1.sleep
		time.Sleep(log.bufTime)
		//2.check file
		log.checkFile()
		//3.check bufsize
		if log.buf.Len() < log.bufSize {
			continue
		}
		//4.write
		log.mutex.Lock()
		if WriteBytes(log.FilePath, log.buf.Bytes()) {
			log.buf.Reset()
		}
		log.mutex.Unlock()
	}

}

func (log *FileLogger) writeAndPrint(title, format string, a ...interface{}) {
	_, fileName, lineNum, _ := runtime.Caller(2)
	stmp := time.Now().Format("2006-01-02 15:04:05")
	var line string
	if title == LogTitleDebug {
		line = fmt.Sprintf("[%-5s] [%s] [%s:%d] [%s]\n", title, stmp, fileName, lineNum, fmt.Sprintf(format, a...))
	} else {
		line = fmt.Sprintf("[%-5s] [%s] [%s]\n", title, stmp, fmt.Sprintf(format, a...))
	}
	if log.console {
		fmt.Print(line)
	}
	log.mutex.Lock()
	log.buf.WriteString(line)
	log.mutex.Unlock()

}

func (log *FileLogger) onExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	<-c
	log.checkFile()
	WriteBytes(log.FilePath, log.buf.Bytes())
	os.Exit(1)
}

func (log *FileLogger) checkFile() {
	size := FileSize(log.FilePath)
	if size < 0 {
		panic(errors.New("open log file failed"))
	}
	if size > log.fileMaxSize {
		stmp := time.Now().Format("2006_01_02_15_04_05_")
		if RenameFile(log.FilePath, stmp+log.FilePath) {
			OpenFile(log.FilePath).Close()
		}
	}

}
