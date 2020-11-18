package log

import "os"

type AbstractLogger struct {
	Logger
}

func (l *AbstractLogger) Debug(msg string, context ...interface{}) {
	l.Log(DEBUG, msg, context...)
}
func (l *AbstractLogger) Info(msg string, context ...interface{}) {
	l.Log(INFO, msg, context...)
}
func (l *AbstractLogger) Warn(msg string, context ...interface{}) {
	l.Log(WARN, msg, context...)
}
func (l *AbstractLogger) Error(msg string, context ...interface{}) {
	l.Log(ERROR, msg, context...)
}
func (l *AbstractLogger) Panic(msg string, context ...interface{}) {
	l.Log(PANIC, msg, context...)
	panic(msg)
}
func (l *AbstractLogger) Fatal(msg string, context ...interface{}) {
	l.Log(FATAL, msg, context...)
	os.Exit(1)
}
