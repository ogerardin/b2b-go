package log4go

import (
	"github.com/sirupsen/logrus"
)

// A Category is a node in a hierarchical Logger tree
type Category struct {
	Name       string
	Priority   logrus.Level
	Appenders  []Appender
	Additivity bool
	//private
	parent *Category
	FieldLogger
}

var _ SimpleLogger = (*Category)(nil)

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

func (c *Category) prepare() {
	effectivePriority := c.getEffectivePriority()
	debugf("  effective priority for %s: %s", c.Name, effectivePriority)

	effectiveAppenders := c.getEffectiveAppenders()
	debugf("  effective appenders for %s: %s", c.Name, effectiveAppenders)

	ca := NewCompositeAppender(effectiveAppenders)

	c.FieldLogger = NewAppenderLogger(effectivePriority, ca)
}

func (c *Category) SetPriority(level logrus.Level) {
	c.Priority = level
}

func (c *Category) AddAppender(appender Appender) {
	c.Appenders = append(c.Appenders, appender)
}
