package log4go

import (
	"github.com/sirupsen/logrus"
)

type Category struct {
	Name       string
	Priority   logrus.Level
	Appenders  []Appender
	Additivity bool
	//private
	parent            *Category
	effectivePriority logrus.Level
	CompositeLogger
}

func (c *Category) getEffectivePriority() logrus.Level {
	if c.Priority != UndefinedLevel {
		return c.Priority
	}
	if c.parent == nil {
		//should not happen since root logger has defined priority
		return logrus.DebugLevel
	}
	return c.parent.getEffectivePriority()
}

func (c *Category) getEffectiveAppenders() []Appender {
	appenders := c.Appenders
	if !c.Additivity {
		return appenders
	}
	if c.parent == nil {
		return appenders
	}
	parentAppenders := c.parent.getEffectiveAppenders()
	appenders = append(appenders, parentAppenders...)
	return appenders
}
