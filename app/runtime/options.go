package runtime

import (
	"b2b-go/app"
	"flag"
	"log"
	"os"
	"runtime"
)

type Options struct {
	confDir        string
	resetDatabase  bool
	resetDeltaIdxs bool
	ShowVersion    bool
	showPaths      bool
	showDeviceId   bool
	doUpgrade      bool
	doUpgradeCheck bool
	upgradeTo      string
	noBrowser      bool
	browserOnly    bool
	HideConsole    bool
	logFile        string
	auditEnabled   bool
	auditFile      string
	verbose        bool
	paused         bool
	unpaused       bool
	guiAddress     string
	guiAPIKey      string
	generateDir    string
	noRestart      bool
	profiler       string
	assetDir       string
	CpuProfile     bool
	stRestarting   bool
	logFlags       int
	ShowHelp       bool
}

func ParseCommandLineOptions() Options {
	options := defaultRuntimeOptions()

	flag.StringVar(&options.confDir, "home", "", "Set configuration directory")
	flag.BoolVar(&options.ShowVersion, "version", false, "Show version")
	flag.BoolVar(&options.ShowHelp, "help", false, "Show this help")
	flag.BoolVar(&options.showDeviceId, "device-id", false, "Show the device ID")
	flag.BoolVar(&options.verbose, "verbose", false, "Print verbose log output")
	if runtime.GOOS == "windows" {
		// Allow user to hide the console window
		flag.BoolVar(&options.HideConsole, "no-console", false, "Hide console window")
	}

	flag.Usage = app.UsageFor(flag.CommandLine, "b2b [options]", "")
	flag.Parse()

	if len(flag.Args()) > 0 {
		flag.Usage()
		os.Exit(2)
	}

	return options
}

func defaultRuntimeOptions() Options {
	options := Options{
		logFlags: log.Ltime,
	}

	if runtime.GOOS != "windows" {
		// On non-Windows, we explicitly default to "-" which means stdout. On
		// Windows, the blank options.logFile will later be replaced with the
		// default path, unless the user has manually specified "-" or
		// something else.
		options.logFile = "-"
	}

	return options
}
