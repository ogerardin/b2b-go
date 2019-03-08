package runtime

import (
	slavemongo "b2b-go/lib/slave-mongo"
	"flag"
	"fmt"
	"github.com/globalsign/mgo"
	"go.uber.org/dig"
	"io/ioutil"
	"os"
)

// application-wide context
var Container *dig.Container

func init() {
	if flag.Lookup("test.v") != nil {
		fmt.Println("Running under 'go test'")
	}
}

func InitMainContainer(test bool) {
	Container = dig.New()

	//Container.Provide(func() Options {
	//	return parseCommandLineOptions()
	//})

	var mgoPath = "mongo-data"
	if test {
		mgoPath, _ = ioutil.TempDir(os.TempDir(), "mongo-test")
	}

	Container.Provide(func() *slavemongo.DBServer {
		server := slavemongo.DBServer{}
		server.SetPath(mgoPath)
		//server.SetPort(27017)
		return &server
	})

	Container.Provide(func(dbs *slavemongo.DBServer) *mgo.Session {
		session := dbs.Session()
		return session
	})
}
