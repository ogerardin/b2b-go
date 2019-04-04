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
	"github.com/containous/staert"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"os"
	"runtime/pprof"
	"strings"
)

var (
	rootCommand = &flaeg.Command{
		Name:                  "b2b",
		Description:           "Peer-to-peer backup",
		Config:                &runtime.CurrentConfig,
		DefaultPointersConfig: &runtime.DefaultPointersConfig,
	}
)

func main() {
	f := flaeg.New(rootCommand, os.Args[1:])

	_, err := f.Parse(rootCommand)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(-1)
	}

	if runtime.CurrentConfig.Version {
		printVersion()
		os.Exit(0)
	}

	// load configuration according to active profiles
	s := staert.NewStaert(rootCommand)

	profiles := strings.Split(runtime.CurrentConfig.Profiles, ",")
	profiles = util.Map(profiles, strings.TrimSpace)
	fmt.Printf("Active profiles: %v\n", profiles)

	s.AddSource(staert.NewTomlSource("b2b", []string{"./conf", "."}))
	for _, profile := range profiles {
		s.AddSource(staert.NewTomlSource("b2b-"+profile, []string{"./conf", "."}))
	}
	s.AddSource(f)

	_, err = s.LoadConfig()
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(-1)
	}

	fmt.Printf("%+v\n", runtime.CurrentConfig)

	startDaemon(&runtime.CurrentConfig)

}

func printVersion() error {
	fmt.Printf("%s %s\n", meta.Version, meta.GitHash)
	return nil
}

func providers(conf *runtime.Configuration) fx.Option {
	providers := fx.Options()

	providers = fx.Options(providers,
		fx.Provide(runtime.DBServerProvider),
	)

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
