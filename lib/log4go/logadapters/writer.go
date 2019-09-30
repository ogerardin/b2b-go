package logadapters

import (
	"b2b-go/lib/log4go"
	"fmt"
	"github.com/sirupsen/logrus"
)

type WriterAdapter struct {
	Logger log4go.LevelLogger
	Level  logrus.Level
}

func (wa *WriterAdapter) Write(p []byte) (n int, err error) {
	wa.Logger.Log(wa.Level, string(p))
	return len(p), nil
}

type PrinterAdapter struct {
	Logger log4go.LevelLogger
	Level  logrus.Level
}

func (pa *PrinterAdapter) Printf(f string, v ...interface{}) {
	pa.Logger.Log(pa.Level, fmt.Sprint(v))
}
