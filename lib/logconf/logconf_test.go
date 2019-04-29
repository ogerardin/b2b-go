package logconf

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func Test(t *testing.T) {

	console := NewConsoleAppender()
	file := NewFileAppender("test.log")

	config = &Config{
		Loggers: []LogAppender{
			{
				Name:     "",
				Level:    logrus.InfoLevel,
				Appender: console,
			},
			{
				Name:     "a",
				Level:    logrus.InfoLevel,
				Appender: file,
			},
			{
				Name:     "a/b",
				Level:    logrus.DebugLevel,
				Appender: console,
			},
		},
	}

	logger := NewLogger("a/b")

	logger.Print("bla")

}
