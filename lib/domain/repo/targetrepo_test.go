package repo

import (
	"b2b-go/lib/domain"
	"b2b-go/lib/runtime"
	"b2b-go/lib/slave-mongo"
	"fmt"
	"github.com/globalsign/mgo"
	"log"
	"testing"
)

func TestTargetRepo(t *testing.T) {

	err := runtime.Container.Invoke(func(dbs *slave_mongo.DBServer) {
		session := dbs.Session()
		defer session.Close()

		testTargetRepoWithSession(t, session)
	})
	if err != nil {
		log.Panicf("Failed to invoke test: %v", err)
	}

}

func testTargetRepoWithSession(t *testing.T, session *mgo.Session) {
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
	id1 := saveTarget(t, repo, target)
	id2 := saveTarget(t, repo, target2)
	loadTarget(t, repo, id1)
	loadTarget(t, repo, id2)
	//time.Sleep(time.Hour)
}

func saveTarget(t *testing.T, repo TargetRepo, target domain.BackupTarget) interface{} {
	id, err := repo.SaveNew(target)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(id)
	return id
}

func loadTarget(t *testing.T, repo TargetRepo, id interface{}) {
	loaded, err := repo.GetById(id)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("loaded %s", (loaded).Desc())
}
