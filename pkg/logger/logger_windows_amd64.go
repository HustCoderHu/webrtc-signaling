package logger

const (
	LOG_EMERG = iota
	LOG_ERR
	LOG_WARNING
	LOG_NOTICE
	LOG_INFO
	LOG_DEBUG
)

func init() {
	logLevel = LOG_INFO
}

// InitLog 初始化日志到syslog
// 初始化之前的日志会打印到标准输出
func InitLog() error {
	return nil
}
