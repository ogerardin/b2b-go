package log4go

import (
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
func GetPackageLogger() Logger {
	debugf("GetPackageLogger")
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic(errors.New("Failed to access call stack"))
	}

	funcName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(funcName, ".")
	packageName := strings.Join(parts[0:len(parts)-1], ".")
	loggerName := strings.ReplaceAll(packageName, "/", ".")

	return GetLogger(loggerName)
}

func GetLogger(name string) Logger {
	return getConfig().GetLogger(name)
}
