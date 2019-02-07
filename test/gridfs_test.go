package test_test

import (
	"b2b-go/lib/slave-mongo"
	"bytes"
	"github.com/globalsign/mgo"
	"io/ioutil"
	"os"
	"testing"
)

var (
	server  slave_mongo.DBServer
	session *mgo.Session
)

func TestGridFs(t *testing.T) {

	var err error

	defer cleanup()

	d, _ := ioutil.TempDir(os.TempDir(), "mongotools-test")

	server = slave_mongo.DBServer{}
	server.SetPath(d)

	session = server.Session()
	db := session.DB("testgridfs")

	fs := db.GridFS("bucket0")
	f, err := fs.Create("sample.txt")
	if err != nil {
		t.Fatal(err)
	}

	fileBytes, err := ioutil.ReadFile("sample.txt")
	if err != nil {
		t.Fatal(err)
	}

	_, err = f.Write(fileBytes)
	if err != nil {
		t.Fatal(err)
	}
	_ = f.Close()

	f, err = fs.Open("sample.txt")
	if err != nil {
		t.Fatal(err)
	}

	storedBytes := make([]byte, f.Size())
	_, err = f.Read(storedBytes)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(fileBytes, storedBytes) != 0 {
		t.Fatal("retrieved bytes differ from file bytes")
	}
}

func cleanup() {
	if session != nil {
		session.Close()
	}
	server.Wipe()
	server.Stop()
}
