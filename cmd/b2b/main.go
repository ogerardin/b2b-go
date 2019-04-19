package main

import (
	"b2b-go/app/repo"
	"b2b-go/app/rest"
	"b2b-go/app/runtime"
	"b2b-go/lib/util"
	"b2b-go/meta"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
)

func main() {

	var conf = runtime.DefaultConfig()

	// parse command line options
	parseCommandLine(conf)

	// if version requested, do it and exit
	if conf.Version {
		printVersion()
		os.Exit(0)
	}

	// load configuration according to active profiles
	err := loadExternalConfig(conf)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(-1)
	}

	//fmt.Printf("%+v\n", conf)
	confBytes, err := json.MarshalIndent(conf, " ", "  ")
	log.Infof("Conf %s", string(confBytes))

	// start the thing
	startApp(conf)

}

func parseCommandLine(conf *runtime.Configuration) {

	pflag.Bool("version", false, "print version information and quit")
	pflag.String("profiles", "", "comma separated list of active profiles")
	pflag.Int("port", conf.Port, "listening port")
	//TODO other flags
	pflag.Parse()

	v := viper.New()
	err := v.BindPFlags(pflag.CommandLine)
	if err != nil {
		panic(errors.Wrap(err, "failed to parse command line"))
	}

	err = v.Unmarshal(conf)
	if err != nil {
		panic(errors.Wrap(err, "failed to unmarshal command line"))
	}

	conf.Viper = v
}

func loadExternalConfig(conf *runtime.Configuration) error {
	profiles := strings.Split(conf.Profiles, ",")
	profiles = util.Map(profiles, strings.TrimSpace)
	fmt.Printf("Active profiles: %v\n", profiles)

	v := conf.Viper
	v.AddConfigPath(".")
	v.SetConfigName("b2b")
	err := conf.Viper.MergeInConfig()
	if err != nil {
		log.Warn(errors.Wrap(err, "Failed to load main config"))
	}

	for _, profile := range profiles {
		v.SetConfigName("b2b-" + profile)
		err := v.MergeInConfig()
		if err != nil {
			log.Warn(errors.Wrapf(err, "Failed to merge profile %s", profile))
		}
	}

	err = v.Unmarshal(conf)
	if err != nil {
		panic(errors.Wrap(err, "Failed to unmarshall config"))
	}
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

func startGin(lc fx.Lifecycle, g *gin.Engine, conf *runtime.Configuration) {

	lc.Append(fx.Hook{
		OnStart: func(c context.Context) error {
			go g.Run(":" + strconv.Itoa(conf.Port))
			return nil
		},
		//OnStop: func(c context.Context) error {
		//	log.Print("Stopping")
		//	return nil
		//},
	})

}
