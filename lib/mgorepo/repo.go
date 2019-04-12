package mgorepo

import (
	"b2b-go/lib/typeregistry"
	"b2b-go/lib/util"
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

func New(s *mgo.Session, coll string) *Repo {
	repo := Repo{
		session: s,
		coll:    coll,
	}
	return &repo
}

func (r *Repo) SaveNew(item interface{}) (bson.ObjectId, error) {
	session := r.session.Copy()
	defer session.Close()

	coll := session.DB("").C(r.coll)

	id := bson.NewObjectId()
	wrapper := wrap(item, id)
	err := coll.Insert(wrapper)

	return id, err
}

func (r *Repo) Update(id bson.ObjectId, item interface{}) error {
	session := r.session.Copy()
	defer session.Close()

	coll := session.DB("").C(r.coll)

	wrapper := wrap(item, id)
	err := coll.UpdateId(id, wrapper)

	return err
}

func (r *Repo) GetById(id bson.ObjectId) (interface{}, error) {
	session := r.session.Copy()
	defer session.Close()

	coll := session.DB("").C(r.coll)

	// read wrapper
	var w map[string]interface{}
	err := coll.Find(bson.M{"_id": id}).One(&w)
	if err != nil {
		return nil, err
	}

	i, err := unwrap(w)
	if err != nil {
		return nil, errors.Wrap(err, "Unwrapping generated an error")
	}

	return i, nil
}

func (r *Repo) GetAll(result interface{}) error {

	resultv := reflect.ValueOf(result)
	if resultv.Kind() != reflect.Ptr {
		log.Panic("result argument must be a slice address")
	}

	slicev := resultv.Elem()

	if slicev.Kind() == reflect.Interface {
		slicev = slicev.Elem()
	}
	if slicev.Kind() != reflect.Slice {
		log.Panic("result argument must be a slice address")
	}

	session := r.session.Copy()
	defer session.Close()

	coll := session.DB("").C(r.coll)

	iter := coll.Find(bson.M{}).Iter()

	var errs = make([]string, 0)
	var w map[string]interface{}
	for iter.Next(&w) {
		item, err := unwrap(w)
		if err != nil {
			errs = append(errs, err.Error())
		} else {
			slicev = reflect.Append(slicev, reflect.ValueOf(item))
		}
	}
	err := iter.Close()
	if err != nil {
		return err
	}

	if len(errs) > 0 {
		return errors.New("Unwrapping generated errors: " + strings.Join(errs, ", "))
	}

	resultv.Elem().Set(slicev)
	return nil
}

func (r *Repo) Delete(id bson.ObjectId) error {
	session := r.session.Copy()
	defer session.Close()

	coll := session.DB("").C(r.coll)

	return coll.RemoveId(id)
}

func wrap(item interface{}, id bson.ObjectId) map[string]interface{} {
	value := util.ConcreteValue(item)
	t := value.Type()

	if !(t.Kind() == reflect.Struct) {
		log.Panic("expected struct")
	}

	itemAsMap := structToMap(value.Interface())

	itemAsMap["_id"] = id
	itemAsMap["_t"] = typeregistry.GetKey(t)

	return itemAsMap
}

func structToMap(item interface{}) map[string]interface{} {

	var result map[string]interface{}
	err := mapstructure.Decode(item, &result)
	if err != nil {
		log.Panic(err)
	}

	return result
}

func unwrap(item map[string]interface{}) (interface{}, error) {

	// Obtain Type from from its key from the type registry
	// The type must have been previously registered with typeregistry.Register
	typeKey := item["_t"].(string)
	t := typeregistry.GetType(typeKey)
	if t == nil {
		log.Panicf("Unknown type key '%s' - did you forget to register the type?", typeKey)
	}

	// get a pointer to a new value of this type
	pt := reflect.New(t)

	// populate value from wrapper.V
	err := mapstructure.Decode(item, pt.Interface())
	if err != nil {
		return nil, err
	}

	// return the value as interface{}
	i := pt.Elem().Interface()
	return i, nil

}
