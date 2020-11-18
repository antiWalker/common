package wrapper

import (
	"github.com/sirupsen/logrus"
	"github.com/antiWalker/common/gsr/log"
	"os"
)

type logger struct {
	log.Logger
	Logrus *logrus.Logger
}

func NewLogger() *logger {
	l := &logger{}
	l.Logrus = logrus.New()
	l.Logrus.SetFormatter(&logrus.JSONFormatter{})
	l.Logrus.SetLevel(logrus.TraceLevel)
	return l
}

func (l *logger) Debug(msg string, context ...interface{}) {
	l.withContext(context...).Debug(msg)
}

func (l *logger) Info(msg string, context ...interface{}) {
	l.withContext(context...).Info(msg)
}

func (l *logger) Warn(msg string, context ...interface{}) {
	l.withContext(context...).Warn(msg)
}

func (l *logger) Error(msg string, context ...interface{}) {
	l.withContext(context...).Error(msg)
}

func (l *logger) Panic(msg string, context ...interface{}) {
	l.withContext(context...).Panic(msg)
	panic(msg)
}

func (l *logger) Fatal(msg string, context ...interface{}) {
	l.withContext(context...).Fatal(msg)
	os.Exit(1)
}

func (l *logger) Log(level string, msg string, context ...interface{}) {
	parseLevel, _ := logrus.ParseLevel(level)
	if len(context) > 0 {
		l.Logrus.WithField("context", context[0]).Log(parseLevel, msg)
	} else {
		l.Logrus.Log(parseLevel, msg)
	}
}

func (l *logger) withContext(context ...interface{}) logrus.FieldLogger {
	if len(context) > 0 {
		return l.Logrus.WithField("context", context[0])
	} else {
		return l.Logrus
	}
}
