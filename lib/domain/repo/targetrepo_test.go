package repo

import (
	"b2b-go/lib/domain"
	"b2b-go/lib/slave-mongo"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestTargetRepo(t *testing.T) {
	d, _ := ioutil.TempDir(os.TempDir(), "mongotools-test")
	server := slave_mongo.DBServer{}
	server.SetPath(d)
	//server.SetPort(27017)
	session := server.Session()
	defer server.Stop()
	defer session.Close()

	repo := NewTargetRepo(session)

	target := domain.LocalTarget{
		BackupTargetBase: domain.BackupTargetBase{
			Name: "target 1",
		},
	}
	target2 := domain.PeerTarget{
		BackupTargetBase: domain.BackupTargetBase{
			Name: "target 2",
		},
		Hostname: "peerhost",
		Port:     9999,
	}

	id1 := save(t, repo, target)
	id2 := save(t, repo, target2)

	load(t, repo, id1)
	load(t, repo, id2)

	//time.Sleep(time.Hour)

}

func save(t *testing.T, repo TargetRepo, target domain.BackupTarget) interface{} {
	id1, err := repo.SaveNew(target)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(id1)
	return id1
}

func load(t *testing.T, repo TargetRepo, id interface{}) {
	loaded, err := repo.GetById(id)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("loaded %s", (loaded).Desc())
}
