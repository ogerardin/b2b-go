package main

import (
	"b2b-go/app/repo"
	"b2b-go/app/rest"
	"b2b-go/app/runtime"
	"b2b-go/lib/util"
	"b2b-go/meta"
	"context"
	"fmt"
	"github.com/containous/flaeg"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"os"
	"runtime/pprof"
)

var (
	rootCommand = &flaeg.Command{
		Name:                  "b2b",
		Description:           "Peer-to-peer backup",
		Config:                &runtime.CurrentConfig,
		DefaultPointersConfig: &runtime.DefaultPointersConfig,
	}

	runCommand = &flaeg.Command{
		Name:                  "run",
		Description:           "Start the daemon and web UI server",
		Config:                &runtime.CurrentConfig,
		DefaultPointersConfig: &runtime.DefaultPointersConfig,
		Run: func() error {
			return startDaemon(&runtime.CurrentConfig)
		},
	}

	versionCommand = &flaeg.Command{
		Name:                  "version",
		Description:           "print version information and exits",
		Config:                &runtime.CurrentConfig,
		DefaultPointersConfig: &runtime.DefaultPointersConfig,
		Run: func() error {
			return printVersion()
		},
	}
)

func main() {
	f := flaegConfig(os.Args[1:])

	rootCommand.Run = func() error {
		// When no command is specified, the root command is invoked. In this case we want to just print the help
		// message; unfortunately flaeg has no easy way to do it.
		// This hack recreates a flaeg config identical to the main one, and simulates a call with "--help"
		f := flaegConfig([]string{"", "--help"})
		f.Run()
		return nil
	}

	err := f.Run()

	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
}

func flaegConfig(args []string) *flaeg.Flaeg {
	f := flaeg.New(rootCommand, args)
	f.AddCommand(runCommand)
	f.AddCommand(versionCommand)
	return f
}

func printVersion() error {
	fmt.Printf("%s %s", meta.Version, meta.GitHash)
	os.Exit(0)
	return nil
}

func startDaemon(conf *runtime.Configuration) error {

	app := fx.New(
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

func handleOptions(lc fx.Lifecycle, conf runtime.Configuration) error {
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
