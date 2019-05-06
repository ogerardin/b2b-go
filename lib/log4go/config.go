package log4go

import (
	"github.com/sirupsen/logrus"
	"io"
	"math"
	"strings"
)

const UndefinedLevel = math.MaxUint32

type Appender struct {
	name      string
	Formatter logrus.Formatter
	Writer    io.Writer
}

func (a *Appender) String() string {
	return a.name
}

type LogAppender struct {
	Name     string
	Appender *Appender
	Level    logrus.Level
}

type Config struct {
	RootLogger *Category
	Loggers    map[string]*Category
}

func DefaultConfig() *Config {
	return &Config{
		RootLogger: &Category{
			Name:   "ROOT",
			parent: nil,
			Appenders: []Appender{
				NewConsoleAppender(),
			},
			Priority: logrus.DebugLevel,
		},
		Loggers: nil,
	}
}

func loadConfig() *Config {
	//TODO load from external config!
	return DefaultConfig()
}

func (conf *Config) GetLogger(name string) Logger {
	debugf("GetLogger('%s')", name)
	logger, found := conf.Loggers[name]
	if found {
		debugf("  returning logger from cache")
		return logger
	}

	category := &Category{
		Name:       name,
		Additivity: true,
		Priority:   UndefinedLevel, //inherited
		Appenders:  nil,
	}

	conf.AddLogger(category)

	return category
}

func (conf *Config) AddLogger(category *Category) {
	parent := conf.getParent(category.Name)
	category.parent = parent
	category.effectivePriority = UndefinedLevel
	conf.insertLogger(category)
	conf.updateEffectivePriorities()
}

func (conf *Config) getParent(name string) *Category {
	parent := conf.RootLogger

	for name, logger := range conf.Loggers {
		if strings.HasPrefix(name, logger.Name) && len(logger.Name) > len(parent.Name) {
			parent = logger
		}
	}
	debugf("  parent for '%s' is %v ", name, parent)
	return parent
}

func (conf *Config) insertLogger(category *Category) {
	for _, logger := range conf.Loggers {
		if logger.parent == category.parent {
			logger.parent = category
		}
	}

	conf.Loggers[category.Name] = category
}

func (conf *Config) updateEffectivePriorities() {
	for _, logger := range conf.Loggers {
		logger.effectivePriority = logger.getEffectivePriority()
	}
}
