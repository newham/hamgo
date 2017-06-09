package hamgo

import "testing"
import "time"

func Test_Debug(t *testing.T) {
	logFile := "./test.log"
	deleteFile(logFile)
	newLogger(logFile)
	for i := 0; i < 10000; i++ {
		Log.Debug("test:%d", i)
		time.Sleep(10)
	}
	time.Sleep(1000)
}

func Test_setFormat(t *testing.T) {

}
