package main

import (
	"b2b-go/app/repo"
	"b2b-go/app/rest"
	"b2b-go/app/runtime"
	"b2b-go/meta"
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"log"
	"os"
	"runtime/pprof"
)

func main() {

	options := runtime.ParseCommandLineOptions()
	handleUsage(options)

	app := fx.New(
		fx.Provide(func() runtime.Options { return options }),
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

func handleUsage(options runtime.Options) {
	if options.ShowVersion {
		fmt.Printf("%s %s", meta.Version, meta.GitHash)
		os.Exit(0)
	}

	if options.ShowHelp {
		flag.Usage()
		os.Exit(0)
	}

}

func handleOptions(lc fx.Lifecycle, options runtime.Options) error {
	if options.HideConsole {
		//osutil.HideConsole()
	}

	lc.Append(fx.Hook{
		OnStop: func(c context.Context) error {
			if options.CpuProfile {
				pprof.StopCPUProfile()
			}
			return nil
		},
	})

	return nil
}

func startGin(lc fx.Lifecycle, g *gin.Engine) {

	lc.Append(fx.Hook{
		OnStart: func(c context.Context) error {
			go g.Run(":8080")
			return nil
		},
		OnStop: func(c context.Context) error {
			log.Print("Stopping")
			return nil
		},
	})

}
