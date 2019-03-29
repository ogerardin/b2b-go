package util

import "github.com/sirupsen/logrus"

func RecoverPanicAndLog(log *logrus.Logger, msg string) {
	if r := recover(); r != nil {
		log.Debugf(msg+": %s", r)
	}
}
