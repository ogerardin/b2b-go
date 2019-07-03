package log4go

import "github.com/sirupsen/logrus"

// a logger where the log level is passed as a parameter
type LevelLogger interface {
	Log(level logrus.Level, args ...interface{})
	Logf(level logrus.Level, fmt string, args ...interface{})
	Logln(level logrus.Level, args ...interface{})
}

// An adapter for FieldLogger to implement LevelLogger
type LevelLoggerFieldLogger struct {
	logrus.FieldLogger
}

func (cl LevelLoggerFieldLogger) Log(level logrus.Level, args ...interface{}) {
	switch level {
	case logrus.PanicLevel:
		cl.Panic(args...)
	case logrus.FatalLevel:
		cl.Fatal(args...)
	case logrus.ErrorLevel:
		cl.Error(args...)
	case logrus.WarnLevel:
		cl.Warn(args...)
	case logrus.InfoLevel:
		cl.Info(args...)
	case logrus.DebugLevel:
		cl.Debug(args...)
	}
}

func (cl LevelLoggerFieldLogger) Logf(level logrus.Level, fmt string, args ...interface{}) {
	switch level {
	case logrus.PanicLevel:
		cl.Panicf(fmt, args...)
	case logrus.FatalLevel:
		cl.Fatalf(fmt, args...)
	case logrus.ErrorLevel:
		cl.Errorf(fmt, args...)
	case logrus.WarnLevel:
		cl.Warnf(fmt, args...)
	case logrus.InfoLevel:
		cl.Infof(fmt, args...)
	case logrus.DebugLevel:
		cl.Debugf(fmt, args...)
	}
}

func (cl LevelLoggerFieldLogger) Logln(level logrus.Level, args ...interface{}) {
	switch level {
	case logrus.PanicLevel:
		cl.Panicln(args...)
	case logrus.FatalLevel:
		cl.Fatalln(args...)
	case logrus.ErrorLevel:
		cl.Errorln(args...)
	case logrus.WarnLevel:
		cl.Warnln(args...)
	case logrus.InfoLevel:
		cl.Infoln(args...)
	case logrus.DebugLevel:
		cl.Debugln(args...)
	}
}
