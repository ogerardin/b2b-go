package log4go

import "github.com/sirupsen/logrus"

type LevelLogger interface {
	Log(level logrus.Level, args ...interface{})
	Logf(level logrus.Level, fmt string, args ...interface{})
	Logln(level logrus.Level, args ...interface{})
}

type WriterAdapter struct {
	Logger LevelLogger
	Level  logrus.Level
}

func (wa *WriterAdapter) Write(p []byte) (n int, err error) {
	wa.Logger.Log(wa.Level, string(p))
	return len(p), nil
}
