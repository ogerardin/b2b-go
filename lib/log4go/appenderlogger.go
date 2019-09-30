package log4go

import (
	"b2b-go/lib/util"
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

var _ Logger = (*appenderLoggerChainable)(nil)

func (al *AppenderLogger) maybeAppend(level logrus.Level, msg string) {
	if level >= al.level {
		al.appender.Append(level, logrus.Fields{}, msg)
	}
}

func (al *AppenderLogger) Debug(msg string) {
	al.maybeAppend(logrus.DebugLevel, msg)
}

func (al *AppenderLogger) Info(msg string) {
	al.maybeAppend(logrus.InfoLevel, msg)
}

func (al *AppenderLogger) Warn(msg string) {
	al.maybeAppend(logrus.WarnLevel, msg)
}

func (al *AppenderLogger) Fatal(msg string) {
	al.maybeAppend(logrus.FatalLevel, msg)
}

func (al *AppenderLogger) Error(msg string) {
	al.maybeAppend(logrus.ErrorLevel, msg)
}

func (al *AppenderLogger) Log(level logrus.Level, msg string) {
	al.maybeAppend(level, msg)
}

func (al *AppenderLogger) WithField(key string, value interface{}) FieldLogger {
	return al.WithFields(logrus.Fields{key: value})
}

func (al *AppenderLogger) WithError(err error) FieldLogger {
	return al.WithField(logrus.ErrorKey, err)
}

func (al *AppenderLogger) WithFields(fields logrus.Fields) FieldLogger {
	chainable := &appenderLoggerChainable{
		AppenderLogger: al,
		fields:         fields,
	}
	return chainable
}

func (chainable *appenderLoggerChainable) WithFields(fields logrus.Fields) FieldLogger {
	util.MergeMaps(chainable.fields, fields)
	return chainable
}

func (chainable *appenderLoggerChainable) maybeAppend(level logrus.Level, msg string) {
	if level < chainable.level {
		return
	}
	chainable.appender.Append(level, chainable.fields, msg)
}

func NewAppenderLogger(level logrus.Level, appender Appender) *AppenderLogger {
	return &AppenderLogger{level: level, appender: appender}
}
