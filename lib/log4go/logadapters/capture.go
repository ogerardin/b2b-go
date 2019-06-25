package logadapters

import (
	"b2b-go/lib/log4go"
	"bufio"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

func CaptureStdOut() *os.File {
	return capture(&os.Stdout)
}

func CaptureStdErr() *os.File {
	return capture(&os.Stderr)
}

func capture(pf **os.File) *os.File {
	r, w, err := os.Pipe()
	if err != nil {
		panic(errors.Wrap(err, "Failed to create pipe"))
	}
	*pf = w

	return r
}

func Feed(r io.Reader, log log4go.LevelLogger, level logrus.Level) {
	pump(r, &WriterAdapter{
		Logger: log,
		Level:  level,
	})

}

func pump(r io.Reader, w io.Writer) {
	go func() {
		reader := bufio.NewReader(r)
		writer := bufio.NewWriter(w)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				panic(errors.Wrap(err, "ReadString returned error"))
			}
			line = strings.TrimSuffix(line, "\n")
			if len(line) == 0 {
				continue
			}
			writer.WriteString(line)
			writer.Flush()
		}
	}()

	return
}
