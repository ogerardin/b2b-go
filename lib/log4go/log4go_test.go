package log4go

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func Test(t *testing.T) {

	testConfig()

	logger := GetLogger("a.b.c")

	logger.Debug("bla")

	logger.Warn("this is a warning: WARNING!")

	logger.WithField("color", "blue").Info("We have a color")

}

func testConfig() {
	debug = true

	console := NewConsoleAppender()
	file := NewFileAppender("test.log")

	config = DefaultConfig()
	config.getRootLogger().SetPriority(logrus.InfoLevel)

	config.AddNode(&Category{
		Name:       "a",
		Priority:   logrus.InfoLevel,
		Appenders:  []Appender{file},
		Additivity: false,
	},
	)
	config.AddNode(&Category{
		Name:       "a.b",
		Priority:   logrus.DebugLevel,
		Appenders:  []Appender{console},
		Additivity: true,
	},
	)
}

func TestPackage(t *testing.T) {
	testConfig()

	logger := GetPackageLogger()

	logger.Info("Hello logger")
}
