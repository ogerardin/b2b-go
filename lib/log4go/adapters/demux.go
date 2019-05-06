package adapters

import (
	"b2b-go/lib/log4go"
	"github.com/sirupsen/logrus"
	"strings"
)

type DemuxingAdapter struct {
	Outputs []DemuxOutput
}

type DemuxOutput struct {
	Predicate func(string) bool
	Logger    log4go.LevelLogger
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
		if output.Predicate(line) {
			level, msg := output.Extracter(line)
			output.Logger.Log(level, msg)
			return
		}
	}
	// no predicate matches -> no logging
}
