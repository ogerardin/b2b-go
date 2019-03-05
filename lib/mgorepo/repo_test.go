package mgorepo

import (
	slave_mongo "b2b-go/lib/slave-mongo"
	"b2b-go/lib/typeregistry"
	"github.com/globalsign/mgo"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

type I interface {
	String() string
}

type A struct {
	field1 int
}

func (A) String() string {
	return "i'm an A!"
}

type B struct {
	A
	field2 int
}

func (B) String() string {
	return "i'm a B!"
}

type C struct {
	A
	field3 int
}

func (C) String() string {
	return "i'm a C!"
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

func (r *TestRepo) SaveNew(item I) (interface{}, error) {
	saved, err := r.Repo.SaveNew(item)
	return saved, err
}

func (r *TestRepo) GetById(id interface{}) (I, error) {
	retrieved, err := r.Repo.GetById(id)
	return retrieved.(I), err
}

func TestGenericRepo(t *testing.T) {
	d, _ := ioutil.TempDir(os.TempDir(), "mongotools-test")
	server := slave_mongo.DBServer{}
	server.SetPath(d)
	//server.SetPort(27017)
	session := server.Session()
	defer server.Wipe()
	defer server.Stop()
	defer session.Close()

	repo := NewTestRepo(session)

	instanceB := B{
		A: A{
			field1: 1,
		},
	}
	id1, err := repo.SaveNew(&instanceB)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id1)

	instanceC := C{
		A: A{
			field1: 2,
		},
	}
	id2, err := repo.SaveNew(&instanceC)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id2)

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

	//time.Sleep(time.Hour)

}
