package runtime

import (
	slavemongo "b2b-go/lib/slave-mongo"
	"context"
	"flag"
	"github.com/globalsign/mgo"
	"go.uber.org/fx"
	"io/ioutil"
	"log"
	"os"
)

func init() {
	if flag.Lookup("test.v") != nil {
		log.Println("Running under 'go test'")
	}
}

func TestDBServerProvider(lc fx.Lifecycle) *slavemongo.DBServer {
	return dbServer(lc, true)
}

func DBServerProvider(lc fx.Lifecycle) *slavemongo.DBServer {
	return dbServer(lc, false)
}

func dbServer(lc fx.Lifecycle, test bool) *slavemongo.DBServer {
	server := slavemongo.DBServer{}

	if test {
		mgoPath, _ := ioutil.TempDir(os.TempDir(), "mongo-test")
		server.SetPath(mgoPath)
		server.SetPort(27017)
	}

	lc.Append(fx.Hook{
		OnStart: func(c context.Context) error {
			log.Print("Starting slave Mongo server")
			server.Start()
			return nil
		},
		OnStop: func(c context.Context) error {
			log.Print("Stopping slave Mongo server")
			if test {
				server.Wipe()
			}
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
