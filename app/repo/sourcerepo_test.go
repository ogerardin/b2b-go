package repo

import (
	"b2b-go/app/domain"
	"b2b-go/app/runtime"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"testing"
)

func TestSourceRepo(t *testing.T) {
	testApp := fxtest.New(t,
		fx.Provide(func() *testing.T { return t }),
		fx.Provide(runtime.TestDBServerProvider()),
		fx.Provide(runtime.SessionProvider()),

		fx.Invoke(testSourceRepoWithSession),
	)
	testApp.RequireStart()
	testApp.RequireStop()
}

func testSourceRepoWithSession(t *testing.T, session *mgo.Session) {
	repo := NewSourceRepo(session)
	source := domain.FilesystemSource{
		BackupSourceBase: domain.BackupSourceBase{
			Name: "source 1",
		},
		Paths: []string{"temp1"},
	}

	source2 := domain.FilesystemSource{
		BackupSourceBase: domain.BackupSourceBase{
			Name: "source 2",
		},
		Paths: []string{"temp2", "temp3"},
	}

	source3 := domain.PeerSource{
		BackupSourceBase: domain.BackupSourceBase{
			Name: "source 3",
		},
	}

	id1 := saveSource(t, repo, source)
	id2 := saveSource(t, repo, source2)
	id3 := saveSource(t, repo, source3)

	sources, err := repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Found %d sources", len(sources))
	for _, s := range sources {
		t.Log(s.Desc())
	}

	loadSource(t, repo, id1)
	loadSource(t, repo, id2)
	loadSource(t, repo, id3)

	//time.Sleep(time.Hour)
}

func saveSource(t *testing.T, repo SourceRepo, source domain.BackupSource) bson.ObjectId {
	id, err := repo.SaveNew(source)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(id)
	return id
}

func loadSource(t *testing.T, repo SourceRepo, id bson.ObjectId) {
	loadedSource1, err := repo.GetById(id)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("loaded %s", (loadedSource1).Desc())
}
