package log4go

import (
	"github.com/sirupsen/logrus"
	"io"
	"strings"
)

type Appender struct {
	name      string
	Formatter logrus.Formatter
	Writer    io.Writer
}

func (a Appender) String() string {
	return a.name
}

type LogAppender struct {
	Name     string
	Appender *Appender
	Level    logrus.Level
}

type Config struct {
	//Appenders []Appender
	Loggers []LogAppender
}

var (
	config *Config
)

func SetConfig(c *Config) {
	config = c
}

func getConfig() *Config {
	if config != nil {
		return config
	}
	config = loadConfig()
	return config
}

func defaultConfig() *Config {
	return &Config{
		//Appenders: []Appender{},
		Loggers: []LogAppender{},
	}
}

func loadConfig() *Config {
	//TODO load from external config!
	return defaultConfig()
}

func defaultContext() loggerContext {
	return loggerContext{
		appenders: []LogAppender{},
	}
}

func (conf *Config) getContext(name string) loggerContext {
	result := defaultContext()
	result.name = name

	for _, l := range conf.Loggers {
		if strings.HasPrefix(name, l.Name) {
			logAppender := LogAppender{
				Name:     name,
				Level:    l.Level,
				Appender: l.Appender,
			}
			result.appenders = append(result.appenders, logAppender)
			debugf("  logger name matches '%s' -> added destination %s ", l.Name, logAppender)
		}
	}

	return result
}
