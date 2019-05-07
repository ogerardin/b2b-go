package log4go

import (
	"github.com/pkg/errors"
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
			Name:   "",
			parent: nil,
			Appenders: []*Appender{
				NewConsoleAppender(),
			},
			Priority: logrus.DebugLevel,
		},
		Loggers: make(map[string]*Category),
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
		debugf("  returning existing logger")
		if logger.CompositeLogger == nil {
			logger.prepare()
		}
		return logger
	}

	category := &Category{
		Name:       name,
		Additivity: true,
		Priority:   UndefinedLevel, //inherited
		Appenders:  nil,
	}

	conf.AddLogger(category)

	category.prepare()

	return category
}

func (conf *Config) AddLogger(category *Category) {
	debugf("Adding logger: %+v", category)

	if len(category.Name) == 0 {
		panic(errors.New("name cannot be empty"))
	}

	parent := conf.getParent(category.Name)
	debugf("  parent for '%s' is '%s' ", category.Name, parent.Name)
	category.parent = parent

	conf.insertLogger(category)
}

func (conf *Config) getParent(name string) *Category {
	parent := conf.RootLogger
	for _, logger := range conf.Loggers {
		if strings.HasPrefix(name, logger.Name) && len(logger.Name) > len(parent.Name) {
			parent = logger
		}
	}
	return parent
}

func (conf *Config) insertLogger(category *Category) {
	// reassign the parent
	for _, logger := range conf.Loggers {
		if logger.parent == category.parent {
			logger.parent = category
		}
	}
	// store in cache
	conf.Loggers[category.Name] = category
}

func (conf *Config) getRootLogger() *Category {
	return conf.RootLogger
}
