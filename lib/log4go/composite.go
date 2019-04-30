package log4go

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type CompositeLogger struct {
	loggers []logrus.FieldLogger
}

func (cl *CompositeLogger) String() string {
	return fmt.Sprintf("CompositeLogger{loggers:%#v}", cl.loggers)
}

func NewCompositeLogger(loggers ...logrus.FieldLogger) *CompositeLogger {
	return &CompositeLogger{
		loggers: loggers,
	}
}

func (cl *CompositeLogger) WithField(key string, value interface{}) *CompositeLogger {
	entries := make([]logrus.FieldLogger, 0)
	for _, l := range cl.loggers {
		entry := l.WithField(key, value)
		entries = append(entries, logrus.FieldLogger(entry))
	}

	return NewCompositeLogger(entries...)
}

func (cl *CompositeLogger) WithFields(fields logrus.Fields) *CompositeLogger {
	entries := make([]logrus.FieldLogger, 0)
	for _, l := range cl.loggers {
		entry := l.WithFields(fields)
		entries = append(entries, logrus.FieldLogger(entry))
	}

	return NewCompositeLogger(entries...)
}

func (cl *CompositeLogger) WithError(err error) *CompositeLogger {
	entries := make([]logrus.FieldLogger, 0)
	for _, l := range cl.loggers {
		entry := l.WithError(err)
		entries = append(entries, logrus.FieldLogger(entry))
	}

	return NewCompositeLogger(entries...)
}

func (cl *CompositeLogger) Debugf(format string, args ...interface{}) {
	for _, l := range cl.loggers {
		l.Debugf(format, args)
	}
}

func (cl *CompositeLogger) Infof(format string, args ...interface{}) {
	for _, l := range cl.loggers {
		l.Infof(format, args)
	}
}

func (cl *CompositeLogger) Warnf(format string, args ...interface{}) {
	for _, l := range cl.loggers {
		l.Warnf(format, args)
	}
}

func (cl *CompositeLogger) Warningf(format string, args ...interface{}) {
	for _, l := range cl.loggers {
		l.Warningf(format, args)
	}
}

func (cl *CompositeLogger) Errorf(format string, args ...interface{}) {
	for _, l := range cl.loggers {
		l.Errorf(format, args)
	}
}

func (cl *CompositeLogger) Debug(args ...interface{}) {
	for _, l := range cl.loggers {
		l.Debug(args)
	}
}

func (cl *CompositeLogger) Info(args ...interface{}) {
	for _, l := range cl.loggers {
		l.Info(args)
	}
}

func (cl *CompositeLogger) Warn(args ...interface{}) {
	for _, l := range cl.loggers {
		l.Warn(args)
	}
}

func (cl *CompositeLogger) Warning(args ...interface{}) {
	for _, l := range cl.loggers {
		l.Warning(args)
	}
}

func (cl *CompositeLogger) Error(args ...interface{}) {
	for _, l := range cl.loggers {
		l.Error(args)
	}
}

func (cl *CompositeLogger) Debugln(args ...interface{}) {
	for _, l := range cl.loggers {
		l.Debugln(args)
	}
}

func (cl *CompositeLogger) Infoln(args ...interface{}) {
	for _, l := range cl.loggers {
		l.Infoln(args)
	}
}

func (cl *CompositeLogger) Warnln(args ...interface{}) {
	for _, l := range cl.loggers {
		l.Warnln(args)
	}
}

func (cl *CompositeLogger) Warningln(args ...interface{}) {
	for _, l := range cl.loggers {
		l.Warningln(args)
	}
}

func (cl *CompositeLogger) Errorln(args ...interface{}) {
	for _, l := range cl.loggers {
		l.Errorln(args)
	}
}

func (cl *CompositeLogger) Append(logger logrus.FieldLogger) {
	cl.loggers = append(cl.loggers, logger)
}

func (cl CompositeLogger) Print(args ...interface{}) {
	for _, l := range cl.loggers {
		l.Print(args)
	}
}

func (cl CompositeLogger) Printf(fmt string, args ...interface{}) {
	for _, l := range cl.loggers {
		l.Printf(fmt, args)
	}
}

func (cl CompositeLogger) Println(args ...interface{}) {
	for _, l := range cl.loggers {
		l.Println(args)
	}
}

func (cl CompositeLogger) Fatal(args ...interface{}) {
	for _, l := range cl.loggers {
		l.Fatal(args)
	}
}

func (cl CompositeLogger) Fatalf(fmt string, args ...interface{}) {
	for _, l := range cl.loggers {
		l.Fatalf(fmt, args)
	}
}

func (cl CompositeLogger) Fatalln(args ...interface{}) {
	for _, l := range cl.loggers {
		l.Fatalln(args)
	}
}

func (cl CompositeLogger) Panic(args ...interface{}) {
	for _, l := range cl.loggers {
		l.Panic(args)
	}
}

func (cl CompositeLogger) Panicf(fmt string, args ...interface{}) {
	for _, l := range cl.loggers {
		l.Panicf(fmt, args)
	}
}

func (cl CompositeLogger) Panicln(args ...interface{}) {
	for _, l := range cl.loggers {
		l.Panicln(args)
	}
}
