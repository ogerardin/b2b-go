package log4go

import (
	"b2b-go/lib/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
)

// An Appender is just a target for logging messages that writes, stores or forwards a log entry.
// It does not typically filter on the log level.
type Appender interface {
	Append(level logrus.Level, fields logrus.Fields, msg string)
}

// an appender that uses an underlying logrus.FieldLogger to format and write log entries to
// a io.Writer
type LoggerAppender struct {
	logrus.FieldLogger
}

func (l *LoggerAppender) Append(level logrus.Level, fields logrus.Fields, msg string) {
	l.FieldLogger.WithFields(fields).Logln(level, msg)
}

// Returns a newly instantiated LoggerAppender that writes log entries to the console.
func NewConsoleAppender() Appender {
	return &LoggerAppender{
		FieldLogger: &logrus.Logger{
			Out: os.Stdout,
			Formatter: &prefixed.TextFormatter{
				ForceColors:     true,
				ForceFormatting: true,
				FullTimestamp:   true,
			},
			Level: logrus.TraceLevel,
		},
	}
}

func NewConsoleAppender_obsolete() *Appender_obsolete {
	return &Appender_obsolete{
		name: "Console",
		Formatter: &prefixed.TextFormatter{
			ForceColors:     true,
			ForceFormatting: true,
			FullTimestamp:   true,
		},
		Writer: os.Stdout,
	}
}

func NewFileAppender(filename string) Appender {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, util.OS_ALL_RW)
	if err != nil {
		panic(errors.Wrapf(err, "Failed to open file %s for writing", filename))
	}
	return &LoggerAppender{
		FieldLogger: &logrus.Logger{
			Out:       file,
			Formatter: &logrus.TextFormatter{},
			Level:     logrus.TraceLevel,
		},
	}
}

func NewFileAppender_obsolete(filename string) *Appender_obsolete {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, util.OS_ALL_RW)
	if err != nil {
		panic(errors.Wrapf(err, "Failed to open file %s for writing", filename))
	}
	return &Appender_obsolete{
		name:      "File: " + filename,
		Formatter: &logrus.TextFormatter{},
		Writer:    file,
	}
}
