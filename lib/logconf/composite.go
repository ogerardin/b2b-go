package logconf

import "github.com/sirupsen/logrus"

type CompositeLogger struct {
	loggers []logrus.StdLogger
}

func NewCompositeLogger(loggers ...logrus.StdLogger) CompositeLogger {
	return CompositeLogger{
		loggers: loggers,
	}
}

func (cl *CompositeLogger) Append(logger logrus.StdLogger) {
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
