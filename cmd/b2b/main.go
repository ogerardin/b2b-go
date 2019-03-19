package main

import (
	app2 "b2b-go/app"
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

	options := app2.ParseCommandLineOptions()
	handleUsage(options)

	app := fx.New(
		fx.Provide(func() app2.Options { return options }),
		fx.Provide(app2.GinProvider),

		fx.Invoke(handleOptions),
		fx.Invoke(app2.RegisterAppRoutes),
		fx.Invoke(startGin),
	)

	app.Run()
}

func handleUsage(options app2.Options) {
	if options.ShowVersion {
		fmt.Printf("%s %s", meta.Version, meta.GitHash)
		os.Exit(0)
	}

	if options.ShowHelp {
		flag.Usage()
		os.Exit(0)
	}

}

func handleOptions(lc fx.Lifecycle, options app2.Options) error {
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
