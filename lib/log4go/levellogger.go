package log4go

import "github.com/sirupsen/logrus"

type LevelLogger interface {
	Log(level logrus.Level, args ...interface{})
	Logf(level logrus.Level, fmt string, args ...interface{})
	Logln(level logrus.Level, args ...interface{})
}
