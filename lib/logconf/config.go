package logconf

import (
	"github.com/sirupsen/logrus"
	"io"
	"strings"
)

type Appender struct {
	Formatter logrus.Formatter
	Writer    io.Writer
}

type Config struct {
	Appenders []Appender
	Loggers   []LogAppender
}

type LogAppender struct {
	Name     string
	Level    logrus.Level
	Appender Appender
}

var (
	config *Config
)

func getConfig() *Config {
	if config != nil {
		return config
	}
	config = loadConfig()
	return config
}

func loadConfig() *Config {
	//TODO load from external config!
	return defaultConfig()
}

func defaultConfig() *Config {
	return &Config{
		Appenders: []Appender{},
		Loggers:   []LogAppender{},
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
		}
	}

	return result
}
