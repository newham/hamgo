package hamgo

var Log Logger

type Logger interface {
	Error(format string, a ...interface{})
	Info(format string, a ...interface{})
	Debug(format string, a ...interface{})
	Warn(format string, a ...interface{})
}

type FileLogger struct {
	FilePath string
}

func newLogger(filePath string) {
	Log = FileLogger{FilePath: filePath}
}

func (log FileLogger) Error(format string, a ...interface{}) {

}

func (log FileLogger) Info(format string, a ...interface{}) {

}

func (log FileLogger) Debug(format string, a ...interface{}) {

}

func (log FileLogger) Warn(format string, a ...interface{}) {

}
