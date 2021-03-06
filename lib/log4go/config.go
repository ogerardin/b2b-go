package log4go

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"math"
	"strings"
)

const UndefinedLevel = math.MaxUint32

type Config struct {
	RootLogger *Category
	Loggers    map[string]*Category
}

// Returns a default configuration for log4go, with just a root logger logging to a console appender
// at DEBUG level.
func DefaultConfig() *Config {
	return &Config{
		RootLogger: &Category{
			Name:   "",
			parent: nil,
			Appenders: []Appender{
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

// Returns a logger by name. If a logger already exists for this name, it is returned.
// Otherwise a new logger is created and configured according to the configuration
func (conf *Config) GetLogger(name string) FieldLogger {
	debugf("GetLogger('%s')", name)
	logger, found := conf.Loggers[name]
	if found {
		debugf("  logger exists")
		if logger.FieldLogger == nil {
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

	conf.AddNode(category)

	category.prepare()

	return category
}

func (conf *Config) AddNode(category *Category) {
	debugf("Adding node: %+v", category)

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
