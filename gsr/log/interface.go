package log

const (
	DEBUG = "debug"
	INFO  = "info"
	WARN  = "warning"
	ERROR = "error"
	PANIC = "panic"
	FATAL = "fatal"
)

type Logger interface {
	Debug(msg string, context ...interface{})
	Info(msg string, context ...interface{})
	Warn(msg string, context ...interface{})
	Error(msg string, context ...interface{})
	Panic(msg string, context ...interface{})
	Fatal(msg string, context ...interface{})
	Log(level string, msg string, context ...interface{})
}
