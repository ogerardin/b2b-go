package log4go

import "github.com/sirupsen/logrus"

// a logger where the log level is passed as a parameter
type LevelLogger interface {
	Log(level logrus.Level, msg string)
}

// An adapter for FieldLogger to implement LevelLogger
type LevelLoggerFieldLogger struct {
	logrus.FieldLogger
}

func (llfl *LevelLoggerFieldLogger) Log(level logrus.Level, msg string) {
	switch level {
	case logrus.PanicLevel:
		llfl.Panic(msg)
	case logrus.FatalLevel:
		llfl.Fatal(msg)
	case logrus.ErrorLevel:
		llfl.Error(msg)
	case logrus.WarnLevel:
		llfl.Warn(msg)
	case logrus.InfoLevel:
		llfl.Info(msg)
	case logrus.DebugLevel:
		llfl.Debug(msg)
	}
}
