package log4go

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func Test(t *testing.T) {

	testConfig()

	logger := GetLogger("a.b.c")

	logger.Debugln("bla")

	logger.Warnf("this is a warning: %s", "WARNING!")

	logger.WithField("color", "blue").Info("We have a color")

}

func testConfig() {
	console := NewConsoleAppender()
	file := NewFileAppender("test.log")

	config = DefaultConfig()
	config.getRootLogger().SetPriority(logrus.InfoLevel)

	config.AddLogger(&Category{
		Name:       "a",
		Priority:   logrus.InfoLevel,
		Appenders:  []*Appender{file},
		Additivity: false,
	},
	)
	config.AddLogger(&Category{
		Name:       "a.b",
		Priority:   logrus.DebugLevel,
		Appenders:  []*Appender{console},
		Additivity: true,
	},
	)
}

func TestPackage(t *testing.T) {
	testConfig()

	logger := GetDefaultLogger()

	logger.Info("Hello logger")
}
