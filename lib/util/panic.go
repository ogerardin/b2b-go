package util

import "github.com/sirupsen/logrus"

func RecoverPanicAndLog(logger *logrus.Logger, msg string) {
	if r := recover(); r != nil {
		if logger != nil {
			logrus.Debugf(msg+": %s", r)
		} else {
			log.Debugf(msg+": %s", r)
		}
	}
}
