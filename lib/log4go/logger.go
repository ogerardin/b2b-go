package log4go

import "github.com/sirupsen/logrus"

// A modified logrus.FieldLogger interface that does not depend on logrus.Entry
type Logger interface {
	WithField(key string, value interface{}) Logger
	WithFields(fields logrus.Fields) Logger
	WithError(err error) Logger

	Debug(msg string)
	Info(msg string)
	Print(msg string)
	Warn(msg string)
	Warning(msg string)
	Error(msg string)
	Fatal(msg string)
	Panic(msg string)

	// extending LevelLogger interface
	Log(level logrus.Level, args ...interface{})
}
