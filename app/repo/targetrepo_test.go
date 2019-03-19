package repo

import (
	"b2b-go/app"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"testing"
)

func TestTargetRepo(t *testing.T) {
	testApp := fxtest.New(t,
		fx.Provide(func() *testing.T { return t }),
		fx.Provide(app.TestDBServerProvider),
		fx.Provide(app.SessionProvider),

		fx.Invoke(testTargetRepoWithSession),
	)
	testApp.RequireStart()
	testApp.RequireStop()
}

func testTargetRepoWithSession(t *testing.T, session *mgo.Session) {
	repo := app.NewTargetRepo(session)
	target := app.LocalTarget{
		BackupTargetBase: app.BackupTargetBase{
			Name: "target 1",
		},
	}
	target2 := app.PeerTarget{
		BackupTargetBase: app.BackupTargetBase{
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

func saveTarget(t *testing.T, repo app.TargetRepo, target app.BackupTarget) bson.ObjectId {
	id, err := repo.SaveNew(target)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(id)
	return id
}

func loadTarget(t *testing.T, repo app.TargetRepo, id bson.ObjectId) {
	loaded, err := repo.GetById(id)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("loaded %s", (loaded).Desc())
}
