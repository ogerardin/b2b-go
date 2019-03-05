package main

import (
	"b2b-go/lib/rest"
	"b2b-go/lib/usage"
	"b2b-go/meta"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
)

const (
	exitSuccess            = 0
	exitError              = 1
	exitNoUpgradeAvailable = 2
	exitRestarting         = 3
	exitUpgrading          = 4
)

var (
	stop = make(chan int)
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

func main() {

	options := parseCommandLineOptions()

	if options.hideConsole {
		//osutil.HideConsole()
	}

	if options.showVersion {
		fmt.Printf("%s %s", meta.Version, meta.GitHash)
		return
	}

	if options.showHelp {
		flag.Usage()
		return
	}

	b2bMain(options)
}

func b2bMain(runtimeOptions RuntimeOptions) {
	setupSignalHandling()

	go rest.StartApi()

	code := <-stop

	//mainService.Stop()

	log.Print("Exiting")

	if runtimeOptions.cpuProfile {
		pprof.StopCPUProfile()
	}

	os.Exit(code)
}

func setupSignalHandling() {
	// Exit cleanly with "restarting" code on SIGHUP.

	restartSign := make(chan os.Signal, 1)
	sigHup := syscall.Signal(1)
	signal.Notify(restartSign, sigHup)
	go func() {
		<-restartSign
		stop <- exitRestarting
	}()

	// Exit with "success" code (no restart) on INT/TERM

	stopSign := make(chan os.Signal, 1)
	sigTerm := syscall.Signal(15)
	signal.Notify(stopSign, os.Interrupt, sigTerm)
	go func() {
		<-stopSign
		stop <- exitSuccess
	}()
}
