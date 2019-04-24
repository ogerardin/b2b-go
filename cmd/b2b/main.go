package main

import (
	"b2b-go/app/repo"
	"b2b-go/app/rest"
	"b2b-go/app/runtime"
	"b2b-go/lib/util"
	"b2b-go/meta"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	"github.com/heetch/confita/backend/flags"
	"go.uber.org/fx"
	"os"
	"runtime/pprof"
	"strings"
)

func main() {

	var conf = runtime.DefaultConfig()

	// parse command line options
	err := parseCommandLine(conf)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(-1)
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
		os.Exit(-1)
	}

	fmt.Printf("%+v\n", conf)

	// start the thing
	startApp(conf)

}

func parseCommandLine(conf *runtime.Configuration) error {
	loader := confita.NewLoader(
		flags.NewBackend(),
	)

	err := loader.Load(context.Background(), conf)
	return err
}

func loadExternalConfig(conf *runtime.Configuration) error {
	profiles := strings.Split(conf.Profiles, ",")
	profiles = util.Map(profiles, strings.TrimSpace)
	fmt.Printf("Active profiles: %v\n", profiles)

	backends := make([]backend.Backend, 0)
	backends = append(backends, env.NewBackend())
	backends = append(backends, file.NewBackend("conf/b2b.toml"))
	for _, profile := range profiles {
		backends = append(backends, file.NewBackend("conf/b2b-"+profile+".toml"))
	}
	backends = append(backends, flags.NewBackend())

	loader := confita.NewLoader(
		backends...,
	)

	err := loader.Load(context.Background(), conf)
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
