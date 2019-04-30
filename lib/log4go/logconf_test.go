package log4go

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func Test(t *testing.T) {

	testConfig()

	logger := GetLogger("a/b")

	logger.Debugln("bla")

	logger.Warnf("this is a warning: %s", "WARNING!")

	logger.WithField("color", "blue").Info("We have a color")

}

func testConfig() {
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
}

func TestPackage(t *testing.T) {
	testConfig()

	logger := GetDefaultLogger()

	logger.Info("Hello logger")
}
