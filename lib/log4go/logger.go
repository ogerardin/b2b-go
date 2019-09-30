package log4go

import "github.com/sirupsen/logrus"

// Basic logger interface.
// Please note that the formatting variants (*f) are not included as they complexify the interface without real benefit.
type SimpleLogger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

// A modified logrus.FieldLogger interface that does not depend on logrus.Entry
type FieldLogger interface {
	SimpleLogger
	WithField(key string, value interface{}) FieldLogger
	WithFields(fields logrus.Fields) FieldLogger
	WithError(err error) FieldLogger
}

// a logger where the log level is passed as a parameter
type LevelLogger interface {
	Log(level logrus.Level, msg string)
}

type Logger interface {
	FieldLogger
	LevelLogger
}

// adapter for SimpleLogger into LevelLogger
type LevelLoggerAdapter struct {
	SimpleLogger
}

func (l *LevelLoggerAdapter) Log(level logrus.Level, msg string) {
	switch level {
	case logrus.DebugLevel:
		l.Debug(msg)
	case logrus.InfoLevel:
		l.Info(msg)
	case logrus.WarnLevel:
		l.Warn(msg)
	case logrus.ErrorLevel:
		l.Error(msg)
	case logrus.FatalLevel:
		l.Fatal(msg)
	}
}
