package log

var std Logger

func Debug(msg string, context ...interface{}) {
	std.Debug(msg, context...)
}
func Info(msg string, context ...interface{}) {
	std.Info(msg, context...)
}
func Warn(msg string, context ...interface{}) {
	std.Warn(msg, context...)
}
func Error(msg string, context ...interface{}) {
	std.Error(msg, context...)
}
func Panic(msg string, context ...interface{}) {
	std.Panic(msg, context...)
}
func Fatal(msg string, context ...interface{}) {
	std.Fatal(msg, context...)
}
func Log(level string, msg string, context ...interface{}) {
	std.Log(level, msg, context...)
}

func SetLogger(logger Logger) {
	std = logger
}
