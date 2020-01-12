package log4go

import (
	"b2b-go/lib/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"io"
	"os"
)

// An Appender is just a target for logging messages that writes, stores or forwards a log entry.
// It does not typically filter on the log level.
type Appender interface {
	Append(level logrus.Level, fields logrus.Fields, msg string)
}

// an appender that uses an underlying logrus.FieldLogger to format and write log entries to
// a io.Writer
type LogrusAppender struct {
	logrus.FieldLogger
}

func (l *LogrusAppender) Append(level logrus.Level, fields logrus.Fields, msg string) {
	l.FieldLogger.WithFields(fields).Logln(level, msg)
}

func NewLogrusAppender(w io.Writer, formatter logrus.Formatter) *LogrusAppender {
	return &LogrusAppender{
		FieldLogger: &logrus.Logger{
			Out:       w,
			Formatter: formatter,
			// we set the level to TraceLevel so that entries are never filtered
			Level: logrus.TraceLevel,
		},
	}
}

// Returns a newly instantiated Appender that writes log entries to the console.
func NewConsoleAppender() Appender {
	formatter := &prefixed.TextFormatter{
		ForceColors:     true,
		ForceFormatting: true,
		FullTimestamp:   true,
	}
	return NewLogrusAppender(os.Stdout, formatter)
}

// Returns a newly instantiated Appender that writes log entries to the specified file
func NewFileAppender(filename string) Appender {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, util.OS_ALL_RW)
	if err != nil {
		panic(errors.Wrapf(err, "Failed to open file %s for writing", filename))
	}
	formatter := &logrus.TextFormatter{}
	return NewLogrusAppender(file, formatter)
}
