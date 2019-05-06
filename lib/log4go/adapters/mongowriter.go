package adapters

import (
	"b2b-go/lib/log4go"
	"github.com/sirupsen/logrus"
	"strings"
)

var levelMap map[uint8]logrus.Level

func init() {
	levelMap = map[uint8]logrus.Level{
		'D': logrus.DebugLevel,
		'I': logrus.InfoLevel,
		'W': logrus.WarnLevel,
		'E': logrus.ErrorLevel,
		'F': logrus.FatalLevel,
	}
}

type MongoWriterAdapter struct {
	Logger       log4go.LevelLogger
	DefaultLevel logrus.Level
}

func (mwa *MongoWriterAdapter) Write(p []byte) (n int, err error) {
	lines := strings.Split(string(p), "\n")
	for _, line := range lines {
		if len(line) > 29 {
			//skip timestamp
			line = line[29:]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		levelChar := line[0]
		level, found := levelMap[levelChar]
		if found {
			line = line[2:]
		} else {
			level = mwa.DefaultLevel
		}
		mwa.Logger.Log(level, line)
	}
	return len(p), nil
}
