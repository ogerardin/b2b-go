package log4go

import (
	"github.com/sirupsen/logrus"
)

// implements the Appender interface for a set of Appenders
type CompositeAppender struct {
	appenders []Appender
}

var _ Appender = (*CompositeAppender)(nil)

func (ca *CompositeAppender) Append(level logrus.Level, fields logrus.Fields, message string) {
	for _, a := range ca.appenders {
		a.Append(level, fields, message)
	}
}

func NewCompositeAppender(appenders []Appender) *CompositeAppender {
	return &CompositeAppender{
		appenders: appenders,
	}
}
