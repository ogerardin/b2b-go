package repo

import (
	"b2b-go/lib/domain"
	"b2b-go/lib/typeregistry"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

type BackupSourceRepo struct {
	session *mgo.Session
	coll    string
}

type Wrapper struct {
	Id bson.ObjectId `bson:"_id"`
	T  string        `bson:"_t"`
	V  interface{}
}

const (
	defaultColl = "sources"
)

func NewSourceRepo(s *mgo.Session) *BackupSourceRepo {
	repo := BackupSourceRepo{
		session: s,
		coll:    defaultColl,
	}
	return &repo
}

func (bsr *BackupSourceRepo) SaveNew(source domain.BackupSource) (interface{}, error) {
	session := bsr.session.Copy()
	defer session.Close()

	coll := session.DB("").C(bsr.coll)

	t := reflect.ValueOf(source).Elem().Type()
	wrapper := Wrapper{
		Id: bson.NewObjectId(),
		T:  typeregistry.GetKey(t),
		V:  source,
	}

	err := coll.Insert(wrapper)
	return wrapper.Id, err
}

func (bsr *BackupSourceRepo) GetById(id interface{}) (*domain.BackupSource, error) {
	session := bsr.session.Copy()
	defer session.Close()

	coll := session.DB("").C(bsr.coll)

	// read wrapper
	wrapper := Wrapper{}
	err := coll.Find(bson.M{"_id": id}).One(&wrapper)
	if err != nil {
		return nil, err
	}
	//tv := reflect.TypeOf(wrapper.V)
	//fmt.Println(tv)

	// obtain Type from registry
	t := typeregistry.GetType(wrapper.T)

	// get a pointer to a new value of this type
	pt := reflect.New(t)

	// populate value using wrapper.V
	err = mapstructure.Decode(wrapper.V, pt.Interface())
	if err != nil {
		return nil, err
	}

	// return the value as *BackupSource
	i := pt.Elem().Interface().(domain.BackupSource)
	return &i, err
}
