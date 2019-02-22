package repo

import (
	"b2b-go/lib/domain"
	"b2b-go/lib/slave-mongo"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestSave(t *testing.T) {
	d, _ := ioutil.TempDir(os.TempDir(), "mongotools-test")

	server := slave_mongo.DBServer{}
	server.SetPath(d)
	server.SetPort(27017)

	defer server.Stop()

	session := server.Session()

	defer session.Close()

	repo := NewSourceRepo(session)

	source := domain.FilesystemSource{
		BackupSourceBase: domain.BackupSourceBase{
			Name: "s1",
		},
		Paths: []string{"temp1"},
	}
	id1, err := repo.SaveNew(&source)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(id1)

	source2 := domain.FilesystemSource{
		BackupSourceBase: domain.BackupSourceBase{
			Name: "s2",
		},
		Paths: []string{"temp2"},
	}
	id2, err := repo.SaveNew(&source2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(id2)

	loadedSource1, err := repo.GetById(id1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(loadedSource1)

	loadedSource2, err := repo.GetById(id2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(loadedSource2)

	time.Sleep(time.Hour)

}
