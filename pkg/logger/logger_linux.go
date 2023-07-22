//go:build linux
// +build linux

package logger

import (
	"log/syslog"
)

const (
	LOG_EMERG   int = int(syslog.LOG_EMERG)
	LOG_ERR     int = int(syslog.LOG_ERR)
	LOG_WARNING int = int(syslog.LOG_WARNING)
	LOG_NOTICE  int = int(syslog.LOG_NOTICE)
	LOG_INFO    int = int(syslog.LOG_INFO)
	LOG_DEBUG   int = int(syslog.LOG_DEBUG)
)

// InitLog 初始化日志到syslog
// 初始化之前的日志会打印到标准输出
func InitLog() error {
	var err error
	initOnce.Do(func() {
		var syslogger Ilogger
		syslogger, err = syslog.New(syslog.LOG_INFO|syslog.LOG_USER, syslogName)
		if err != nil {
			return
		}
		logger = syslogger
	})
	return err
}
