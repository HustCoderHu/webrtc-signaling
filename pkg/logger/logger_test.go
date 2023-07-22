package logger

import "testing"

func Test(t *testing.T) {
	InitLog() // linux 上执行这行，log 会写进 syslog，否则就输出到 stdout
	SetLogLevel(LOG_INFO)
	Info("default log level: %d", GetLogLevel())
}
