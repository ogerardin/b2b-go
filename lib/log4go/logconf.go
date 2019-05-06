package log4go

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"runtime"
	"strings"
)

// a cache of existing loggers by name
var loggers = make(map[string]*CompositeLogger, 0)

var debug = os.Getenv("debug_log4go") != ""

type loggerContext struct {
	name      string
	appenders []LogAppender
}

func (context *loggerContext) newLogger() *CompositeLogger {
	if len(context.appenders) == 0 {
		log.Printf("Warning: no appender for context '%s'", context.name)
		return NewCompositeLogger()
	}

	loggers := make([]logrus.FieldLogger, 0)

	for _, la := range context.appenders {
		debugf("  creating logger for %s", la)
		logger := logrus.New()
		logger.SetLevel(la.Level)
		logger.SetFormatter(la.Appender.Formatter)
		logger.SetOutput(la.Appender.Writer)
		loggers = append(loggers, logrus.FieldLogger(logger))
	}

	logger := NewCompositeLogger(loggers...)
	debugf("Returning: %s", logger)
	return logger
}

func debugf(fmt string, args ...interface{}) {
	if debug {
		log.Printf("[log4go] "+fmt, args...)
	}

}

func GetLogger(name string) Logger {
	debugf("GetLogger(%s)", name)
	logger, found := loggers[name]
	if found {
		debugf("Returning logger from cache")
		return logger
	}

	context := getContext(name)

	logger = context.newLogger()
	// store in cache
	loggers[name] = logger

	return logger
}

func GetDefaultLogger() Logger {
	debugf("GetDefaultLogger")
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
