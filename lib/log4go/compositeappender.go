package log4go

import (
	"github.com/sirupsen/logrus"
)

// implements the Apender interface for a set of Appenders
type CompositeAppender struct {
	appenders []Appender
}

var _ Appender = (*CompositeAppender)(nil)

func (c *CompositeAppender) Append(level logrus.Level, fields logrus.Fields, message string) {
	for _, a := range c.appenders {
		a.Append(level, fields, message)
	}
}

func NewCompositeAppender(appenders []Appender) *CompositeAppender {
	return &CompositeAppender{
		appenders: appenders,
	}
}
