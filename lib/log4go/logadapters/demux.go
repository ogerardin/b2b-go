package logadapters

import (
	"b2b-go/lib/log4go"
	"github.com/sirupsen/logrus"
	"strings"
)

// A DemuxingAdapter is a Writer that is able to distribute written lines to several Loggers.
type DemuxingAdapter struct {
	Outputs []DemuxOutput
}

// DemuxOutput is the structure that associates a Logger with a predicate function that decides whether this line should
// be sent to this Logger, and an extracter function that extracts both a logging level to use and the actual log message
// that should be used.
type DemuxOutput struct {
	Logger log4go.LevelLogger
	// Should this line be sent to this logger?
	Predicate func(string) bool
	// What level and message should be logged ?
	Extracter func(string) (logrus.Level, string)
}

func (dma *DemuxingAdapter) Write(p []byte) (n int, err error) {
	lines := strings.Split(string(p), "\n")
	for _, line := range lines {
		dma.logLine(line)
	}
	return len(p), nil
}

func (dma *DemuxingAdapter) logLine(line string) {
	for _, output := range dma.Outputs {
		if output.Predicate == nil || output.Predicate(line) {
			level, msg := logrus.DebugLevel, line
			if output.Extracter != nil {
				level, msg = output.Extracter(line)
			}
			output.Logger.Log(level, msg)
			return
		}
	}
	// no predicate matches -> no logging
}
