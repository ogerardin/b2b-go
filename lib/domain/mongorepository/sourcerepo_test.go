package mongorepository

import (
	"b2b-go/lib/domain"
	"b2b-go/lib/slave-mongo"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestSave(t *testing.T) {
	d, _ := ioutil.TempDir(os.TempDir(), "mongotools-test")
	server := slave_mongo.DBServer{}
	server.SetPath(d)
	server.SetPort(27017)
	session := server.Session()
	defer server.Stop()
	defer session.Close()

	repo := domain.NewSourceRepo(session)

	source := domain.FilesystemSource{
		BackupSourceBase: domain.BackupSourceBase{
			Name: "source 1",
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
			Name: "source 2",
		},
		Paths: []string{"temp2", "temp3"},
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
	t.Logf("loaded %s", (*loadedSource1).Desc())

	loadedSource2, err := repo.GetById(id2)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("loaded %s", (*loadedSource2).Desc())

	//time.Sleep(time.Hour)

}
