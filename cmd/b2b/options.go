package main

import (
	"b2b-go/lib/usage"
	"flag"
	"log"
	"os"
	"runtime"
)

type RuntimeOptions struct {
	confDir        string
	resetDatabase  bool
	resetDeltaIdxs bool
	showVersion    bool
	showPaths      bool
	showDeviceId   bool
	doUpgrade      bool
	doUpgradeCheck bool
	upgradeTo      string
	noBrowser      bool
	browserOnly    bool
	hideConsole    bool
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
	cpuProfile     bool
	stRestarting   bool
	logFlags       int
	showHelp       bool
}

func parseCommandLineOptions() RuntimeOptions {
	options := defaultRuntimeOptions()

	flag.StringVar(&options.confDir, "home", "", "Set configuration directory")
	flag.BoolVar(&options.showVersion, "version", false, "Show version")
	flag.BoolVar(&options.showHelp, "help", false, "Show this help")
	flag.BoolVar(&options.showDeviceId, "device-id", false, "Show the device ID")
	flag.BoolVar(&options.verbose, "verbose", false, "Print verbose log output")
	if runtime.GOOS == "windows" {
		// Allow user to hide the console window
		flag.BoolVar(&options.hideConsole, "no-console", false, "Hide console window")
	}

	flag.Usage = usage.UsageFor(flag.CommandLine, "b2b [options]", "")
	flag.Parse()

	if len(flag.Args()) > 0 {
		flag.Usage()
		os.Exit(2)
	}

	return options
}

func defaultRuntimeOptions() RuntimeOptions {
	options := RuntimeOptions{
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
