package log4go

import (
	"b2b-go/lib/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
)

// TODO This is the future definition of Appender. The current Appender (struct) is not generic enough
type NewAppender interface {
	Append(level logrus.Level, fields logrus.Fields, message string)
}

type LoggerAppender struct {
	Logger
}

func (l *LoggerAppender) Append(level logrus.Level, fields logrus.Fields, message string) {
	l.Logger.WithFields(fields).Logln(level, message)
}

func NewConsoleAppender() *Appender {
	return &Appender{
		name: "Console",
		//Formatter: &logrus.TextFormatter{
		//	ForceColors: true,
		//},
		Formatter: &prefixed.TextFormatter{
			ForceColors:     true,
			ForceFormatting: true,
			FullTimestamp:   true,
		},
		Writer: os.Stdout,
	}
}

func NewFileAppender(filename string) *Appender {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, util.OS_ALL_RW)
	if err != nil {
		panic(errors.Wrapf(err, "Failed to open file %s for writing", filename))
	}
	return &Appender{
		name:      "File: " + filename,
		Formatter: &logrus.TextFormatter{},
		Writer:    file,
	}
}
