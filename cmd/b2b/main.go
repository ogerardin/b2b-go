package main

import (
	"b2b-go/app/repo"
	"b2b-go/app/rest"
	"b2b-go/app/runtime"
	"b2b-go/lib/log4go"
	"b2b-go/lib/util"
	"b2b-go/meta"
	"context"
	"encoding/json"
	"fmt"
	"github.com/containous/flaeg"
	"github.com/containous/staert"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"os"
	"runtime/pprof"
	"strings"
)

var log *log4go.CompositeLogger

func init() {
	console := log4go.NewConsoleAppender()

	config := &log4go.Config{
		Loggers: []log4go.LogAppender{
			{
				Name:     "b2b-go",
				Level:    logrus.InfoLevel,
				Appender: &console,
			},
		},
	}
	log4go.SetConfig(config)

	log = log4go.GetDefaultLogger()
}

func main() {

	var conf = runtime.DefaultConfig()

	// parse command line options
	err := parseCommandLine(conf)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	// if version requested, do it and exit
	if conf.Version {
		printVersion()
		os.Exit(0)
	}

	// load configuration according to active profiles
	err = loadExternalConfig(conf)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(2)
	}

	//fmt.Printf("%+v\n", conf)
	confBytes, err := json.MarshalIndent(conf, " ", "  ")
	log.Debugf("Conf %s", string(confBytes))

	// start the thing
	startApp(conf)

}

func loadExternalConfig(conf *runtime.Configuration) error {
	profiles := strings.Split(conf.Profiles, ",")
	profiles = util.Map(profiles, strings.TrimSpace)
	fmt.Printf("Active profiles: %v\n", profiles)

	s := staert.NewStaert(conf.Command)
	s.AddSource(staert.NewTomlSource("b2b", []string{"./conf", "."}))
	for _, profile := range profiles {
		s.AddSource(staert.NewTomlSource("b2b-"+profile, []string{"./conf", "."}))
	}
	s.AddSource(conf.Flaeg)

	_, err := s.LoadConfig()
	return err
}

func command(conf *runtime.Configuration) *flaeg.Command {
	command := &flaeg.Command{
		Name:                  "b2b",
		Description:           "Peer-to-peer backup",
		Config:                conf,
		DefaultPointersConfig: runtime.DefaultPointersConfig(),
	}
	return command
}

func parseCommandLine(conf *runtime.Configuration) error {
	command := command(conf)

	f := flaeg.New(command, os.Args[1:])
	cmd, err := f.Parse(command)
	// store these in the config for later use by Staert
	conf.Flaeg = f
	conf.Command = cmd

	return err
}

func printVersion() error {
	fmt.Printf("%s %s\n", meta.Version, meta.GitHash)
	return nil
}

func providers(constructors ...interface{}) fx.Option {
	options := fx.Options()

	for _, provider := range constructors {
		options = fx.Options(options, fx.Provide(provider))
	}

	return options
}

func startApp(conf *runtime.Configuration) error {

	app := fx.New(
		fx.Logger(log),

		fx.Provide(func() *runtime.Configuration { return conf }),
		fx.Provide(runtime.DBServerProvider),
		fx.Provide(runtime.SessionProvider),
		fx.Provide(repo.NewSourceRepo),
		fx.Provide(repo.NewTargetRepo),
		fx.Provide(rest.GinProvider),

		fx.Invoke(handleOptions),
		fx.Invoke(rest.RegisterAppRoutes),
		fx.Invoke(rest.RegisterSourceRoutes),
		fx.Invoke(startGin),
	)

	app.Run()

	return nil
}

func handleOptions(lc fx.Lifecycle, conf *runtime.Configuration) error {
	if conf.HideConsole {
		//osutil.HideConsole()
	}

	if conf.CpuProfile {
		lc.Append(fx.Hook{
			OnStart: func(c context.Context) error {
				file, err := os.OpenFile("pprof", os.O_CREATE|os.O_WRONLY, util.OS_USER_RW|util.OS_GROUP_R|util.OS_OTH_R)
				if err != nil {
					panic(err)
				}
				err = pprof.StartCPUProfile(file)
				if err != nil {
					panic(err)
				}
				return nil
			},
			OnStop: func(c context.Context) error {
				pprof.StopCPUProfile()
				return nil
			},
		})
	}

	return nil
}

func startGin(lc fx.Lifecycle, g *gin.Engine) {

	lc.Append(fx.Hook{
		OnStart: func(c context.Context) error {
			go g.Run(":8080")
			return nil
		},
		//OnStop: func(c context.Context) error {
		//	log.Print("Stopping")
		//	return nil
		//},
	})

}
