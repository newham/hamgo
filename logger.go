package hamgo

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"time"
)

//Log : use Log to write logs
var Log logger

type logger interface {
	Error(format string, a ...interface{})
	Info(format string, a ...interface{})
	Debug(format string, a ...interface{})
	Warn(format string, a ...interface{})
}

type fileLogger struct {
	FilePath    string
	fileFolder  string
	fileMaxSize int64
	fileSize    int64
	Format      string
	buf         *bytes.Buffer
	mutex       *sync.Mutex
	bufTime     time.Duration
	bufSize     int
	console     bool
}

const (
	defaultFilePath     = "./app.log"
	defaultWriteBufTime = 1000        //ms
	defaultWriteBufSize = 1 * 1024    //B
	defaultFileMaxSize  = 1024 * 1024 //B
	defaultConsole      = true
	defaultFormat       = "[%Title] [%Time] [%File] %Text"
	logTitleDebug       = "Debug"
	logTitleInfo        = "Info"
	logTitleError       = "Error"
	logTitleWarn        = "Warn"
	confFilePath        = "log_file"
	confFileMaxSize     = "log_file_max_size"
	confBufSize         = "log_buf_size"
	confBufTime         = "log_buf_time"
	confConsole         = "log_console"
	confFormat          = "log_format"
	confFormatTitle     = "%Title"
	confFormatFile      = "%File"
	confFormatTime      = "%Time"
	confFormatText      = "%Text"
)

func newLogger(filePath string) {
	logger := &fileLogger{
		FilePath:    filePath,
		fileFolder:  currentPath(filePath),
		Format:      defaultFormat,
		console:     true,
		buf:         new(bytes.Buffer),
		mutex:       new(sync.Mutex),
		fileMaxSize: defaultFileMaxSize,
		fileSize:    0,
		bufTime:     time.Duration(defaultWriteBufTime),
		bufSize:     defaultWriteBufSize}
	//create a thread to write buf to log file
	go logger.writeBuf()
	//listen exit signal
	go logger.onExit()

	Log = logger
}

func newLoggerByConf() {
	logger := &fileLogger{
		FilePath:    Conf.DefaultString(confFilePath, defaultFilePath),
		fileFolder:  currentPath(Conf.DefaultString(confFilePath, defaultFilePath)),
		Format:      Conf.DefaultString(confFormat, defaultFormat),
		console:     Conf.DefaultBool(confConsole, defaultConsole),
		buf:         new(bytes.Buffer),
		mutex:       new(sync.Mutex),
		fileMaxSize: Conf.DefaultInt64(confFileMaxSize, defaultFileMaxSize) * 1024,
		fileSize:    0,
		bufTime:     time.Duration(Conf.DefaultInt64(confBufTime, defaultWriteBufTime)),
		bufSize:     Conf.DefaultInt(confBufSize, defaultWriteBufSize) * 1024}
	//create a thread to write buf to log file
	go logger.writeBuf()
	//listen exit signal
	go logger.onExit()

	Log = logger
}

func (log *fileLogger) Error(format string, a ...interface{}) {
	log.writeAndPrint("Error", format, a...)
}

func (log *fileLogger) Info(format string, a ...interface{}) {
	log.writeAndPrint("Info ", format, a...)
}

func (log *fileLogger) Debug(format string, a ...interface{}) {
	log.writeAndPrint("Debug", format, a...)
}

func (log *fileLogger) Warn(format string, a ...interface{}) {
	log.writeAndPrint("Warn ", format, a...)
}

func (log *fileLogger) writeBuf() {
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
		len := int64(log.buf.Len())
		if writeBytes(log.FilePath, log.buf.Bytes()) {
			log.buf.Reset()
			//5.add size
			log.fileSize = log.fileSize + len
		}
		log.mutex.Unlock()

	}

}

func (log *fileLogger) writeAndPrint(title, format string, a ...interface{}) {
	// _, fileName, lineNum, _ := runtime.Caller(2)
	// stmp := time.Now().Format("2006-01-02 15:04:05")
	// var line string
	// line = fmt.Sprintf("[%-5s] [%s] [%s:%d] %s\n", title, stmp, fileName, lineNum, fmt.Sprintf(format, a...))
	line := log.format(title, fmt.Sprintf(format, a...))

	if log.console {
		fmt.Print(line)
	}
	log.mutex.Lock()
	log.buf.WriteString(line)
	log.mutex.Unlock()

}

func (log *fileLogger) onExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	<-c
	writeBytes(log.FilePath, log.buf.Bytes())
	os.Exit(1)
}

func (log *fileLogger) checkFile() {

	if log.fileSize > log.fileMaxSize {
		stmp := time.Now().Format(log.fileFolder + "2006_01_02_15_04_05.log")
		if renameFile(log.FilePath, stmp) {
			openFile(log.FilePath).Close()
			log.fileSize = 0
		}
	}

}

func (log *fileLogger) format(title, text string) string {

	f := log.Format
	if strings.Contains(f, confFormatTitle) {
		f = strings.Replace(f, confFormatTitle, fmt.Sprintf("%s", title), -1)
	}
	if strings.Contains(f, confFormatFile) {
		_, fileName, lineNum, _ := runtime.Caller(3)
		f = strings.Replace(f, confFormatFile, fmt.Sprintf("%s:%d", fileName, lineNum), -1)
	}
	if strings.Contains(f, confFormatTime) {
		stmp := time.Now().Format("2006-01-02 15:04:05")
		f = strings.Replace(f, confFormatTime, stmp, -1)
	}
	if strings.Contains(f, confFormatText) {
		f = strings.Replace(f, confFormatText, text, -1)
	}
	return fmt.Sprintf("%s\n", f)
}
