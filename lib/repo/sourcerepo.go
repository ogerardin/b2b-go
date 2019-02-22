package repo

import (
	"b2b-go/lib/domain"
	"b2b-go/lib/typeregistry"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"reflect"
)

type BackupSourceRepo struct {
	session *mgo.Session
	coll    string
}

type Wrapper struct {
	Id      bson.ObjectId `bson:"_id"`
	TypeKey string
	Val     interface{}
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
		Id:      bson.NewObjectId(),
		TypeKey: typeregistry.GetKey(t),
		Val:     source,
	}

	err := coll.Insert(wrapper)
	return wrapper.Id, err
}

func (bsr *BackupSourceRepo) GetById(id interface{}) (domain.BackupSource, error) {
	session := bsr.session.Copy()
	defer session.Close()

	coll := session.DB("").C(bsr.coll)
	wrapper := Wrapper{}
	err := coll.Find(bson.M{"_id": id}).One(&wrapper)
	if err != nil {
		return nil, err
	}
	tv := reflect.TypeOf(wrapper.Val)
	fmt.Println(tv)

	t := typeregistry.GetType(wrapper.TypeKey)
	pt := reflect.New(t)
	s := pt.Elem().Interface()

	m := wrapper.Val.(bson.M)
	bsonBytes, _ := bson.Marshal(m)
	bson.Unmarshal(bsonBytes, pt)

	return s.(domain.BackupSource), err
}
