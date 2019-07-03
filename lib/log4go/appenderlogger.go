package log4go

import (
	"b2b-go/lib/util"
	"fmt"
	"github.com/sirupsen/logrus"
)

// a Logger that is directly connected to an Appender
type AppenderLogger struct {
	level    logrus.Level
	appender Appender
}

var _ Logger = (*AppenderLogger)(nil)

type appenderLoggerChainable struct {
	*AppenderLogger
	fields logrus.Fields
}

func (al *AppenderLogger) maybeAppendf(level logrus.Level, fmt string, args ...interface{}) {
	if level < al.level {
		return
	}
	al.appender.Appendf(level, logrus.Fields{}, fmt, args...)
}

func (al *AppenderLogger) maybeAppend(level logrus.Level, args ...interface{}) {
	if level < al.level {
		return
	}
	al.appender.Append(level, logrus.Fields{}, args...)
}

func (al *AppenderLogger) Debugf(format string, args ...interface{}) {
	al.maybeAppendf(logrus.DebugLevel, format, args...)
}

func (al *AppenderLogger) Infof(format string, args ...interface{}) {
	al.maybeAppendf(logrus.InfoLevel, format, args...)
}

func (al *AppenderLogger) Printf(format string, args ...interface{}) {
	al.Infof(format, args...)
}

func (al *AppenderLogger) Warnf(format string, args ...interface{}) {
	al.maybeAppendf(logrus.WarnLevel, format, args...)
}

func (al *AppenderLogger) Warningf(format string, args ...interface{}) {
	al.Warnf(format, args...)
}

func (al *AppenderLogger) Errorf(format string, args ...interface{}) {
	al.maybeAppendf(logrus.ErrorLevel, format, args...)
}

func (al *AppenderLogger) Fatalf(format string, args ...interface{}) {
	al.maybeAppendf(logrus.FatalLevel, format, args...)
	//os.Exit(1)
}

func (al *AppenderLogger) Panicf(format string, args ...interface{}) {
	al.maybeAppendf(logrus.PanicLevel, format, args...)
	//panic(fmt.Sprintf(format, args...))
}

func (al *AppenderLogger) Debug(args ...interface{}) {
	al.maybeAppend(logrus.DebugLevel, args...)
}

func (al *AppenderLogger) Info(args ...interface{}) {
	al.maybeAppend(logrus.InfoLevel, args...)
}

func (al *AppenderLogger) Print(args ...interface{}) {
	al.Info(args...)
}

func (al *AppenderLogger) Warn(args ...interface{}) {
	al.maybeAppend(logrus.WarnLevel, args...)
}

func (al *AppenderLogger) Warning(args ...interface{}) {
	al.Warn(args...)
}

func (al *AppenderLogger) Error(args ...interface{}) {
	al.maybeAppend(logrus.ErrorLevel, args...)
}

func (al *AppenderLogger) Fatal(args ...interface{}) {
	al.maybeAppend(logrus.FatalLevel, args...)
	//os.Exit(1)
}

func (al *AppenderLogger) Panic(args ...interface{}) {
	al.maybeAppend(logrus.PanicLevel, args...)
	//panic(fmt.Sprint(args...))
}

func (al *AppenderLogger) Debugln(args ...interface{}) {
	panic("implement me")
}

func (al *AppenderLogger) Infoln(args ...interface{}) {
	panic("implement me")
}

func (al *AppenderLogger) Println(args ...interface{}) {
	panic("implement me")
}

func (al *AppenderLogger) Warnln(args ...interface{}) {
	panic("implement me")
}

func (al *AppenderLogger) Warningln(args ...interface{}) {
	panic("implement me")
}

func (al *AppenderLogger) Errorln(args ...interface{}) {
	panic("implement me")
}

func (al *AppenderLogger) Fatalln(args ...interface{}) {
	panic("implement me")
}

func (al *AppenderLogger) Panicln(args ...interface{}) {
	panic("implement me")
}

func (al *AppenderLogger) Log(level logrus.Level, args ...interface{}) {
	panic("implement me")
}

func (al *AppenderLogger) Logf(level logrus.Level, fmt string, args ...interface{}) {
	panic("implement me")
}

func (al *AppenderLogger) Logln(level logrus.Level, args ...interface{}) {
	panic("implement me")
}

var _ Logger = (*appenderLoggerChainable)(nil)

func (al *AppenderLogger) WithField(key string, value interface{}) Logger {
	return al.WithFields(logrus.Fields{key: value})
}

func (al *AppenderLogger) WithError(err error) Logger {
	return al.WithField(logrus.ErrorKey, err)
}

func (al *AppenderLogger) WithFields(fields logrus.Fields) Logger {
	chainable := &appenderLoggerChainable{
		AppenderLogger: al,
		fields:         fields,
	}
	return chainable
}

func (chainable *appenderLoggerChainable) WithFields(fields logrus.Fields) Logger {
	util.MergeMaps(chainable.fields, fields)
	return chainable
}

func (chainable *appenderLoggerChainable) maybeAppend(level logrus.Level, f string, args ...interface{}) {
	if level < chainable.level {
		return
	}
	chainable.appender.Append(level, chainable.fields, fmt.Sprintf(f, args))
}

func NewAppenderLogger(level logrus.Level, appender Appender) *AppenderLogger {
	return &AppenderLogger{level: level, appender: appender}
}
