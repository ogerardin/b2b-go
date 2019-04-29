package logconf

import (
	"github.com/sirupsen/logrus"
)

type loggerContext struct {
	appenders []LogAppender
	name      string
}

func (context *loggerContext) newLogger() logrus.StdLogger {
	if len(context.appenders) == 0 {
		logrus.Warnf("No Appender for context %s", context.name)
		return NewCompositeLogger()
	}

	loggers := make([]logrus.StdLogger, 0)

	for _, la := range context.appenders {
		logger := logrus.New()
		logger.SetLevel(la.Level)
		logger.SetFormatter(la.Appender.Formatter)
		logger.SetOutput(la.Appender.Writer)
		loggers = append(loggers, logger)
	}

	return NewCompositeLogger(loggers...)
}

func defaultContext() loggerContext {
	return loggerContext{
		appenders: []LogAppender{},
	}
}

func NewLogger(name string) logrus.StdLogger {
	context := getContext(name)
	return context.newLogger()
}

func getContext(name string) loggerContext {
	config := getConfig()
	context := config.getContext(name)
	return context
}
