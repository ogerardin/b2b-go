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
	session     *mgo.Session
	database    string
	coll        string
	idGenerator IdGenerator
}

func New(session *mgo.Session, database string, coll string) *Repo {
	return NewWithIdGenerator(session, database, coll, defaultGenerator())
}

func NewWithIdGenerator(session *mgo.Session, database string, coll string, generator IdGenerator) *Repo {
	return &Repo{
		session:     session,
		database:    database,
		coll:        coll,
		idGenerator: generator,
	}
}

// Save a new item in the repository. The document's ID is obtained by calling NewId() on the idGenerator.
// In case of error, returns nil and the error; otherwise returns the assigned ID and nil.
func (r *Repo) SaveNew(item interface{}) (interface{}, error) {
	session := r.session.Copy()
	defer session.Close()

	coll := session.DB(r.database).C(r.coll)

	id := r.idGenerator.NewId()
	wrapper := wrap(item, id)
	err := coll.Insert(wrapper)

	if err != nil {
		return nil, errors.Wrapf(err, "Mongo failed to save new document %v with id %v", item, id)
	}
	return id, nil

}

// Update the item having the specified ID with the passed data.
// In case of error, returns the error; otherwise returns nil.
func (r *Repo) Update(id interface{}, item interface{}) error {
	session := r.session.Copy()
	defer session.Close()

	coll := session.DB(r.database).C(r.coll)

	wrapper := wrap(item, id)
	err := coll.UpdateId(id, wrapper)
	if err != nil {
		return errors.Wrapf(err, "Mongo failed to update document with id %v", id)
	}

	return nil
}

// Retrieve the item having the specified ID.
// In case of error, returns nil and the error; otherwise returns the item and nil.
func (r *Repo) GetById(id interface{}) (interface{}, error) {
	session := r.session.Copy()
	defer session.Close()

	coll := session.DB(r.database).C(r.coll)

	var w map[string]interface{}
	err := coll.FindId(id).One(&w)
	if err != nil {
		return nil, errors.Wrapf(err, "Mongo failed to retrieve document with id %v", id)
	}

	result, err := unwrap(w)
	if err != nil {
		return nil, errors.Wrap(err, "Unwrapping generated an error")
	}

	return result, nil
}

// Retrieve all the items from the repository into the specified destination, which must be a pointer to a slice of the
// appropriate type.
// In case of error, returns the error; otherwise returns nil.
func (r *Repo) GetAll(result interface{}) error {

	resultv := reflect.ValueOf(result)
	if resultv.Kind() != reflect.Ptr {
		log.Panic("not a pointer")
	}

	slicev := resultv.Elem()

	if slicev.Kind() == reflect.Interface {
		slicev = slicev.Elem()
	}
	if slicev.Kind() != reflect.Slice {
		log.Panic("not a slice")
	}

	session := r.session.Copy()
	defer session.Close()

	coll := session.DB(r.database).C(r.coll)

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
		return errors.Wrap(err, "Error while closing the query")
	}

	if len(errs) > 0 {
		return errors.New("Unwrapping generated errors: " + strings.Join(errs, ", "))
	}

	resultv.Elem().Set(slicev)
	return nil
}

// Delete the item having the specified ID from the collection.
// In case of error, returns the error; otherwise returns nil.
func (r *Repo) Delete(id interface{}) error {
	session := r.session.Copy()
	defer session.Close()

	coll := session.DB(r.database).C(r.coll)

	err := coll.RemoveId(id)
	if err != nil {
		return errors.Wrapf(err, "Mongo failed to delete document with id %v", id)
	}

	return nil
}

// Converts an item into a value that may be passed to Mongo driver's Insert/Update methods.
// For flexibility this function returns a map with the following additional top-level keys:
//
//  - "_id" is the document identifier.
//
//  - "_t" is the original type's identifier as returned by typeregistry.GetKey()
//
func wrap(item interface{}, id interface{}) map[string]interface{} {
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

	// populate value from map
	err := mapstructure.Decode(item, pt.Interface())
	if err != nil {
		return nil, err
	}

	// return the value as interface{}
	i := pt.Elem().Interface()
	return i, nil

}
