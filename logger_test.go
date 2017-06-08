package hamgo

import "testing"
import "time"

func Test_Logt(t *testing.T) {
	logFile := "./test.log"
	newLogger(logFile)
	for i := 0; i < 100000; i++ {
		Log.Debug("test:%d", i)
		time.Sleep(10)
	}
	time.Sleep(1000)
}
