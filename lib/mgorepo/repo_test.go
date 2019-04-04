package mgorepo

import (
	"b2b-go/lib/slave-mongo"
	"b2b-go/lib/typeregistry"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"testing"
)

type I interface {
	String() string
}

type A struct {
	Field1 int
}

func (a A) String() string {
	return "i'm an A! " + strconv.Itoa(a.Field1)
}

type B struct {
	A
	Field2 int
}

func (b B) String() string {
	return "i'm a B! " + strconv.Itoa(b.Field1)
}

type C struct {
	A
	Field3 int
}

func (c C) String() string {
	return "i'm a C! " + strconv.Itoa(c.Field1)
}

func TestMain(m *testing.M) {

	fmt.Println("before tests")
	exitCode := m.Run()
	fmt.Println("after tests")

	os.Exit(exitCode)
}

func init() {
	typeregistry.Register(reflect.TypeOf((*A)(nil)).Elem())
	typeregistry.Register(reflect.TypeOf((*B)(nil)).Elem())
	typeregistry.Register(reflect.TypeOf((*C)(nil)).Elem())
}

type TestRepo struct {
	*Repo
}

func NewTestRepo(s *mgo.Session) *TestRepo {
	return &TestRepo{
		New(s, "test"),
	}
}

func (r *TestRepo) SaveNew(item I) (bson.ObjectId, error) {
	saved, err := r.Repo.SaveNew(item)
	return saved, err
}

func (r *TestRepo) GetById(id bson.ObjectId) (I, error) {
	retrieved, err := r.Repo.GetById(id)
	return retrieved.(I), err
}

func TestGenericRepo(t *testing.T) {
	server := slave_mongo.DBServer{}

	mongoDataPath, _ := ioutil.TempDir(os.TempDir(), "test")
	server.SetPath(mongoDataPath)

	server.Start()
	defer server.Stop()

	session := server.Session()
	defer session.Close()

	testWithSession(t, session)
}

func testWithSession(t *testing.T, session *mgo.Session) {
	repo := NewTestRepo(session)

	instanceB := B{
		A: A{
			Field1: 1,
		},
	}
	id1, err := repo.SaveNew(&instanceB)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id1)

	instanceC := C{
		A: A{
			Field1: 2,
		},
	}

	id2, err := repo.SaveNew(&instanceC)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id2)

	all := make([]I, 0)
	err = repo.GetAll(&all)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Found %d items: %v", len(all), all)

	loaded1, err := repo.GetById(id1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(loaded1.String())

	loaded2, err := repo.GetById(id2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(loaded2.String())

	instanceC2 := C{
		A: A{
			Field1: 333,
		},
	}
	err = repo.Update(id2, instanceC2)
	if err != nil {
		t.Fatal(err)
	}

	loaded22, err := repo.GetById(id2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(loaded22.String())

}
