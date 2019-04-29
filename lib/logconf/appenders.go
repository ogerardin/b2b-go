package logconf

import (
	"b2b-go/lib/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
)

func NewConsoleAppender() Appender {
	return Appender{
		Formatter: &logrus.TextFormatter{},
		Writer:    os.Stdout,
	}
}

func NewFileAppender(filename string) Appender {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND, util.OS_ALL_RW)
	if err != nil {
		panic(errors.Wrapf(err, "Failed to open file %s for writing", filename))
	}
	return Appender{
		Formatter: &logrus.TextFormatter{},
		Writer:    file,
	}
}
