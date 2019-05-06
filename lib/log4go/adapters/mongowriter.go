package adapters

import (
	"b2b-go/lib/log4go"
	"github.com/sirupsen/logrus"
)

type MongoWriterAdapter struct {
	Logger log4go.LevelLogger
}

func (mwa *MongoWriterAdapter) Write(p []byte) (n int, err error) {
	message := string(p)
	message = message[29:]
	mwa.Logger.Log(logrus.InfoLevel, message)
	return len(p), nil
}
