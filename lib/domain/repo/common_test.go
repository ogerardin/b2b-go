package repo

import (
	"b2b-go/lib/runtime"
	"b2b-go/lib/slave-mongo"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	runtime.InitMainContainer(true)

	exitCode := m.Run()

	err := runtime.Container.Invoke(func(dbs *slave_mongo.DBServer) {
		dbs.Wipe()
		dbs.Stop()
	})
	if err != nil {
		log.Printf("Failed to shutdown Mongo server: %v", err)
	}

	os.Exit(exitCode)
}
