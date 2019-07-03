package log4go

import "github.com/sirupsen/logrus"

// Basic logger interface
type Logger interface {
	Debug(msg string)
	Info(msg string)
	Print(msg string)
	Warn(msg string)
	Warning(msg string)
	Error(msg string)
	Fatal(msg string)
	Panic(msg string)
}

// A modified logrus.FieldLogger interface that does not depend on logrus.Entry
type FieldLogger interface {
	Logger
	WithField(key string, value interface{}) Logger
	WithFields(fields logrus.Fields) Logger
	WithError(err error) Logger
}
