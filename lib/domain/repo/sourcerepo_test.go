package repo

import (
	"b2b-go/lib/domain"
	"b2b-go/lib/runtime"
	"fmt"
	"github.com/globalsign/mgo"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"testing"
)

func TestSourceRepo(t *testing.T) {
	testApp := fxtest.New(t,
		fx.Provide(func() *testing.T { return t }),
		fx.Provide(runtime.TestDBServerProvider),
		fx.Provide(runtime.SessionProvider),

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

	id1 := saveSource(t, repo, source)
	id2 := saveSource(t, repo, source2)

	loadSource(t, repo, id1)
	loadSource(t, repo, id2)

	//time.Sleep(time.Hour)
}

func saveSource(t *testing.T, repo SourceRepo, source domain.BackupSource) interface{} {
	id, err := repo.SaveNew(source)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(id)
	return id
}

func loadSource(t *testing.T, repo SourceRepo, id interface{}) {
	loadedSource1, err := repo.GetById(id)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("loaded %s", (loadedSource1).Desc())
}
