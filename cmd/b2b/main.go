package main

import (
	"b2b-go/lib/rest"
	"b2b-go/lib/runtime"
	"b2b-go/meta"
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"log"
	"runtime/pprof"
)

func main() {

	app := fx.New(
		fx.Provide(runtime.OptionsProvider),
		fx.Provide(rest.GinProvider),

		fx.Invoke(handleOptions),
		fx.Invoke(startGin),
	)

	app.Run()
}

func handleOptions(lc fx.Lifecycle, options runtime.Options) error {
	if options.HideConsole {
		//osutil.HideConsole()
	}

	if options.ShowVersion {
		fmt.Printf("%s %s", meta.Version, meta.GitHash)
		//FIXME not an error, just abort app start
		return errors.New("done")
	}

	if options.ShowHelp {
		flag.Usage()
		//FIXME not an error, just abort app start
		return errors.New("done")
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
