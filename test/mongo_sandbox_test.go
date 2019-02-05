package test_test

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"io/ioutil"
	"os"
	"testing"
	"time"

	dbtest "b2b-go/lib/controlled-mongo"
)

type Person struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Name      string
	Phone     string
	Timestamp time.Time
}

func TestEmbedded(t *testing.T) {

	var err error

	d, _ := ioutil.TempDir(os.TempDir(), "mongotools-test")

	server := dbtest.DBServer{}
	server.SetPath(d)

	// Note that the server will be started automagically
	session := server.Session()

	// Insert data programmatically as needed
	c := session.DB("test").C("people")
	err = c.Insert(&Person{Name: "Ale", Phone: "+55 53 1234 4321", Timestamp: time.Now()},
		&Person{Name: "Cla", Phone: "+66 33 1234 5678", Timestamp: time.Now()})

	if err != nil {
		panic(err)
	}

	// Query One
	result := Person{}
	err = c.Find(bson.M{"name": "Ale"}).Select(bson.M{"phone": 0}).One(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println("Phone", result)

	// Query All
	var results []Person
	err = c.Find(bson.M{"name": "Ale"}).Sort("-timestamp").All(&results)

	if err != nil {
		panic(err)
	}
	fmt.Println("Results All: ", results)

	// Update
	colQuerier := bson.M{"name": "Ale"}
	change := bson.M{"$set": bson.M{"phone": "+86 99 8888 7777", "timestamp": time.Now()}}
	err = c.Update(colQuerier, change)
	if err != nil {
		panic(err)
	}

	// Query All
	err = c.Find(bson.M{"name": "Ale"}).Sort("-timestamp").All(&results)

	if err != nil {
		panic(err)
	}
	fmt.Println("Results All: ", results)

	// We can not use "defer session.Close()" because...
	session.Close()

	// ... "server.Wipe()" will panic if there are still connections open
	// for example because you did a .Copy() on the
	// original session in your code. VERY useful!
	server.Wipe()

	// Tear down the server
	server.Stop()
}
