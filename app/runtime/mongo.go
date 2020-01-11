package runtime

import (
	"b2b-go/lib/log4go"
	slavemongo "b2b-go/lib/slave-mongo"
	"b2b-go/lib/util"
	"context"
	"flag"
	"github.com/globalsign/mgo"
	"go.uber.org/fx"
	"io/ioutil"
	"os"
)

var log log4go.FieldLogger

func init() {
	if flag.Lookup("test.v") != nil {
		log.Info("Running under 'go test'")
	}

	log = log4go.GetPackageLogger()
}

func DBServerProvider(lc fx.Lifecycle, conf *Configuration) *slavemongo.DBServer {
	server := slavemongo.DBServer{}

	if conf.MongoDataPath == "" {
		conf.MongoDataPath, _ = ioutil.TempDir(os.TempDir(), "mongo-test")
	}
	os.MkdirAll(conf.MongoDataPath, util.OS_USER_RWX|util.OS_GROUP_RWX|util.OS_OTH_RWX)
	server.SetPath(conf.MongoDataPath)

	server.SetPort(conf.MongoPort)

	/*	server.SetLogAdapter(&logadapters.MongoWriterAdapter{
			Logger:       &log4go.LevelLoggerAdapter{log4go.GetLogger("mongo")},
			DefaultLevel: logrus.DebugLevel,
		})
	*/
	lc.Append(fx.Hook{
		OnStart: func(c context.Context) error {
			log.Info("Starting slave Mongo server")
			server.Start()
			return nil
		},
		OnStop: func(c context.Context) error {
			log.Info("Stopping slave Mongo server")
			//if test {
			//	server.Wipe()
			//}
			server.Stop()
			return nil
		},
	})

	return &server
}

func SessionProvider(lc fx.Lifecycle, server *slavemongo.DBServer) *mgo.Session {
	session := server.Session()

	lc.Append(fx.Hook{
		OnStop: func(c context.Context) error {
			session.Close()
			return nil
		},
	})

	return session
}
