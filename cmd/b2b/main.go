package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	Version           = "unknown-dev"
	Codename          = "Bla Bli"
	BuildStamp        = "0"
	BuildDate         time.Time
	BuildHost         = "unknown"
	BuildUser         = "unknown"
	IsRelease         bool
	IsCandidate       bool
	IsBeta            bool
	LongVersion       string
	BuildTags         []string
	allowedVersionExp = regexp.MustCompile(`^v\d+\.\d+\.\d+(-[a-z0-9]+)*(\.\d+)*(\+\d+-g[0-9a-f]+)?(-[^\s]+)?$`)
)

const (
	exitSuccess            = 0
	exitError              = 1
	exitNoUpgradeAvailable = 2
	exitRestarting         = 3
	exitUpgrading          = 4
)

func setBuildMetadata() {
	// Check for a clean release build. A release is something like
	// "v0.1.2", with an optional suffix of letters and dot separated
	// numbers like "-beta3.47". If there's more stuff, like a plus sign and
	// a commit hash and so on, then it's not a release. If it has a dash in
	// it, it's some sort of beta, release candidate or special build. If it
	// has "-rc." in it, like "v0.14.35-rc.42", then it's a candidate build.
	//
	// So, every build that is not a stable release build has IsBeta = true.
	// This is used to enable some extra debugging (the deadlock detector).
	//
	// Release candidate builds are also "betas" from this point of view and
	// will have that debugging enabled. In addition, some features are
	// forced for release candidates - auto upgrade, and usage reporting.

	exp := regexp.MustCompile(`^v\d+\.\d+\.\d+(-[a-z]+[\d\.]+)?$`)
	IsRelease = exp.MatchString(Version)
	IsCandidate = strings.Contains(Version, "-rc.")
	IsBeta = strings.Contains(Version, "-")

	stamp, _ := strconv.Atoi(BuildStamp)
	BuildDate = time.Unix(int64(stamp), 0)

	date := BuildDate.UTC().Format("2006-01-02 15:04:05 MST")
	LongVersion = fmt.Sprintf(`b2b %s "%s" (%s %s-%s) %s@%s %s`, Version, Codename, runtime.Version(), runtime.GOOS, runtime.GOARCH, BuildUser, BuildHost, date)

	if len(BuildTags) > 0 {
		LongVersion = fmt.Sprintf("%s [%s]", LongVersion, strings.Join(BuildTags, ", "))
	}
}

var (
	stop = make(chan int)
)

const (
	usage = "b2b [options]"
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

	flag.Usage = usageFor(flag.CommandLine, usage, "")
	flag.Parse()

	if len(flag.Args()) > 0 {
		flag.Usage()
		os.Exit(2)
	}

	return options
}

func main() {
	setBuildMetadata()

	options := parseCommandLineOptions()

	if options.hideConsole {
		//osutil.HideConsole()
	}

	if options.showVersion {
		fmt.Println(LongVersion)
		return
	}

	if options.showHelp {
		flag.Usage()
		return
	}

	syncthingMain(options)
}

func syncthingMain(runtimeOptions RuntimeOptions) {
	setupSignalHandling()

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
