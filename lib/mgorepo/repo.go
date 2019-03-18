package mgorepo

import (
	"b2b-go/lib/typeregistry"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"log"
	"reflect"
	"strings"
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

func New(s *mgo.Session, coll string) *Repo {
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
	value := reflect.ValueOf(source)
	var t reflect.Type
	if value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr {
		t = value.Elem().Type()
	} else {
		t = value.Type()
	}

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

	i, err := unwrap(wrapper)
	if err != nil {
		return nil, errors.Wrap(err, "Unwrapping generated an error")
	}

	return i, nil
}

func (r *Repo) GetAll() ([]interface{}, error) {
	session := r.session.Copy()
	defer session.Close()

	coll := session.DB("").C(r.coll)

	// read wrapper
	var wrappers []wrapper
	err := coll.Find(bson.M{}).All(&wrappers)
	if err != nil {
		return nil, err
	}

	var result = make([]interface{}, 0)
	var errs = make([]string, 0)
	for _, w := range wrappers {
		v, err := unwrap(w)
		if err != nil {
			errs = append(errs, err.Error())
		} else {
			result = append(result, v)
		}
	}
	if len(errs) > 0 {
		return nil, errors.New("Unwrapping generated errors: " + strings.Join(errs, " "))
	}

	return result, nil
}

func unwrap(w wrapper) (interface{}, error) {
	// Obtain Type from from its key from the type registry
	// The type must have been previously registered with typeregistry.Register
	t := typeregistry.GetType(w.T)
	if t == nil {
		log.Panicf("Unknown type key '%s' - did you forget to register the type?", w.T)
	}

	// get a pointer to a new value of this type
	pt := reflect.New(t)

	// populate value from wrapper.V
	err := mapstructure.Decode(w.V, pt.Interface())
	if err != nil {
		return nil, err
	}

	// return the value as interface{}
	i := pt.Elem().Interface()
	return i, err

}
