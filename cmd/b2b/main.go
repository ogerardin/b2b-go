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
	runCommand = &flaeg.Command{
		Name:                  "run",
		Config:                &runtime.CurrentConfig,
		DefaultPointersConfig: &runtime.DefaultPointersConfig,
		Run: func() error {
			startDaemon()
			return nil
		},
	}

	versionCommand = &flaeg.Command{
		Name:                  "version",
		Config:                &runtime.CurrentConfig,
		DefaultPointersConfig: &runtime.DefaultPointersConfig,
		Run: func() error {
			printVersion()
			return nil
		},
	}
)

func main() {
	f := flaeg.New(runCommand, os.Args[1:])
	f.AddCommand(versionCommand)

	err := f.Run()

	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
}

func printVersion() {
	fmt.Printf("%s %s", meta.Version, meta.GitHash)
	os.Exit(0)
}

func startDaemon() {

	app := fx.New(
		fx.Provide(func() runtime.Configuration { return runtime.CurrentConfig }),
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
