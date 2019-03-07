package main

import (
	"b2b-go/lib/rest"
	"b2b-go/lib/runtime"
	"b2b-go/meta"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
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

func main() {

	runtime.InitMainContainer(false)

	runtime.Container.Invoke(func(options runtime.Options) {
		b2bMain(options)
	})

}

func b2bMain(options runtime.Options) {
	if options.HideConsole {
		//osutil.HideConsole()
	}

	if options.ShowVersion {
		fmt.Printf("%s %s", meta.Version, meta.GitHash)
		return
	}

	if options.ShowHelp {
		flag.Usage()
		return
	}
	setupSignalHandling()

	go rest.StartApi()

	code := <-stop

	//mainService.Stop()

	log.Print("Exiting")

	if options.CpuProfile {
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
