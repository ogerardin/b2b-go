package log4go

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

var loggers = make(map[string]*CompositeLogger, 0)

type loggerContext struct {
	appenders []LogAppender
	name      string
}

func (context *loggerContext) newLogger() *CompositeLogger {
	if len(context.appenders) == 0 {
		logrus.Warnf("No Appender for context %s", context.name)
		return NewCompositeLogger()
	}

	loggers := make([]logrus.FieldLogger, 0)

	for _, la := range context.appenders {
		logger := logrus.New()
		logger.SetLevel(la.Level)
		logger.SetFormatter(la.Appender.Formatter)
		logger.SetOutput(la.Appender.Writer)
		loggers = append(loggers, logrus.FieldLogger(logger))
	}

	return NewCompositeLogger(loggers...)
}

func defaultContext() loggerContext {
	return loggerContext{
		appenders: []LogAppender{},
	}
}

func GetLogger(name string) *CompositeLogger {
	logger, found := loggers[name]
	if found {
		return logger
	}

	context := getContext(name)
	logger = context.newLogger()
	loggers[name] = logger

	return logger
}

func GetDefaultLogger() *CompositeLogger {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic(errors.New("Failed to access call stack"))
	}

	funcName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(funcName, ".")
	packageName := strings.Join(parts[0:len(parts)-1], ".")

	return GetLogger(packageName)
}

func getContext(name string) loggerContext {
	config := getConfig()
	context := config.getContext(name)
	return context
}
