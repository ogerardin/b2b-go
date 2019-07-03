package log4go

import (
	"github.com/sirupsen/logrus"
)

// A type of logger that forwards log messages to a set of subloggers of type log4rus.FieldLogger.
// Implements both log4go.Logger and log4go.LevelLogger
type CompositeLogger struct {
	loggers []logrus.FieldLogger
}

var (
	// assert *CompositeLogger implements Logger
	_ Logger = (*CompositeLogger)(nil)
	// assert *CompositeLogger implements LevelLogger
	_ LevelLogger = (*CompositeLogger)(nil)
)

func NewCompositeLogger(loggers ...logrus.FieldLogger) *CompositeLogger {
	return &CompositeLogger{
		loggers: loggers,
	}
}

func (cl *CompositeLogger) Log(level logrus.Level, msg string) {
	switch level {
	case logrus.PanicLevel:
		cl.Panic(msg)
	case logrus.FatalLevel:
		cl.Fatal(msg)
	case logrus.ErrorLevel:
		cl.Error(msg)
	case logrus.WarnLevel:
		cl.Warn(msg)
	case logrus.InfoLevel:
		cl.Info(msg)
	case logrus.DebugLevel:
		cl.Debug(msg)
	}
}

/*func (cl *CompositeLogger) String() string {
	if cl == nil {
		return "nil"
	}
	return fmt.Sprintf("CompositeLogger{loggers:%#v}", cl.loggers)
}
*/

func (cl *CompositeLogger) Append(logger logrus.FieldLogger) {
	cl.loggers = append(cl.loggers, logger)
}

func (cl *CompositeLogger) WithField(key string, value interface{}) Logger {
	entries := make([]logrus.FieldLogger, 0)
	for _, l := range cl.loggers {
		entry := l.WithField(key, value)
		entries = append(entries, entry)
	}

	return NewCompositeLogger(entries...)
}

func (cl *CompositeLogger) WithFields(fields logrus.Fields) Logger {
	entries := make([]logrus.FieldLogger, 0)
	for _, l := range cl.loggers {
		entry := l.WithFields(fields)
		entries = append(entries, entry)
	}

	return NewCompositeLogger(entries...)
}

func (cl *CompositeLogger) WithError(err error) Logger {
	entries := make([]logrus.FieldLogger, 0)
	for _, l := range cl.loggers {
		entry := l.WithError(err)
		entries = append(entries, entry)
	}

	return NewCompositeLogger(entries...)
}

func (cl *CompositeLogger) Debug(msg string) {
	for _, l := range cl.loggers {
		l.Debug(msg)
	}
}

func (cl *CompositeLogger) Info(msg string) {
	for _, l := range cl.loggers {
		l.Info(msg)
	}
}

func (cl *CompositeLogger) Warn(msg string) {
	for _, l := range cl.loggers {
		l.Warn(msg)
	}
}

func (cl *CompositeLogger) Warning(msg string) {
	for _, l := range cl.loggers {
		l.Warning(msg)
	}
}

func (cl *CompositeLogger) Error(msg string) {
	for _, l := range cl.loggers {
		l.Error(msg)
	}
}

func (cl CompositeLogger) Print(msg string) {
	for _, l := range cl.loggers {
		l.Print(msg)
	}
}

func (cl CompositeLogger) Fatal(msg string) {
	for _, l := range cl.loggers {
		l.Fatal(msg)
	}
}

func (cl CompositeLogger) Panic(msg string) {
	for _, l := range cl.loggers {
		l.Panic(msg)
	}
}
