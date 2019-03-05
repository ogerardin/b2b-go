package genericrepo

import (
	"b2b-go/lib/typeregistry"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/mapstructure"
	"log"
	"reflect"
)

type Repo struct {
	session *mgo.Session
	coll    string
}

type wrapper struct {
	Id bson.ObjectId `bson:"_id"`
	T  string        `bson:"_t"`
	V  interface{}
}

func NewRepo(s *mgo.Session, coll string) *Repo {
	repo := Repo{
		session: s,
		coll:    coll,
	}
	return &repo
}

func (r *Repo) SaveNew(source interface{}) (interface{}, error) {
	session := r.session.Copy()
	defer session.Close()

	coll := session.DB("").C(r.coll)

	// The value is saved as a wrapper value, with V being the actual value and T being its type key.
	t := reflect.ValueOf(source).Elem().Type()
	wrapper := wrapper{
		Id: bson.NewObjectId(),
		T:  typeregistry.GetKey(t),
		V:  source,
	}

	err := coll.Insert(wrapper)
	return wrapper.Id, err
}

func (r *Repo) GetById(id interface{}) (interface{}, error) {
	session := r.session.Copy()
	defer session.Close()

	coll := session.DB("").C(r.coll)

	// read wrapper
	wrapper := wrapper{}
	err := coll.Find(bson.M{"_id": id}).One(&wrapper)
	if err != nil {
		return nil, err
	}
	//tv := reflect.TypeOf(wrapper.V)
	//fmt.Println(tv)

	// Obtain Type from from its key from the type registry
	// The type must have been previously registered with typeregistry.Register
	t := typeregistry.GetType(wrapper.T)
	if t == nil {
		log.Panicf("Unknown type key '%s' - did you forget to register the type?", wrapper.T)
	}

	// get A pointer to A new value of this type
	pt := reflect.New(t)

	// populate value from wrapper.V
	err = mapstructure.Decode(wrapper.V, pt.Interface())
	if err != nil {
		return nil, err
	}

	// return the value as interface{}
	i := pt.Elem().Interface()
	return i, err
}
