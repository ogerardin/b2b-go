package log4go

import (
	"bufio"
	"github.com/pkg/errors"
	"log"
	"os"
	"runtime"
	"strings"
)

// a cache of existing loggers by name
//var loggers = make(map[string]Logger, 0)

// should we print debug messages?
var debug = os.Getenv("debug_log4go") != ""

var (
	config = DefaultConfig()
)

func SetConfig(c *Config) {
	config = c
}

func getConfig() *Config {
	if config != nil {
		return config
	}
	config = loadConfig()
	return config
}

// internal debug function
func debugf(fmt string, args ...interface{}) {
	if debug {
		log.Printf("[log4go] "+fmt, args...)
	}

}

// returns a Logger with a name based on the current method's package
func GetDefaultLogger() Logger {
	debugf("GetDefaultLogger")
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic(errors.New("Failed to access call stack"))
	}

	funcName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(funcName, ".")
	packageName := strings.Join(parts[0:len(parts)-1], ".")

	return GetLogger(packageName)
}

func GetLogger(name string) Logger {
	return getConfig().GetLogger(name)
}

func CaptureStdOut() {
	stdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		panic(errors.Wrap(err, "Failed to create pipe"))
	}
	os.Stdout = w

	go func() {
		reader := bufio.NewReader(r)
		writer := bufio.NewWriter(stdout)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			writer.WriteString("[stdout] " + line + "\n")
		}

	}()
}

func CaptureStdErr() {
	stderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		panic(errors.Wrap(err, "Failed to create pipe"))
	}
	os.Stderr = w

	go func() {
		reader := bufio.NewReader(r)
		writer := bufio.NewWriter(stderr)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			writer.WriteString("[stderr] " + line + "\n")
		}

	}()
}
