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

type HasId interface {
	GetId() string
	SetId(id string)
}

// this should be a const but compiler complains "const initializer is not a constant"
var hasIdInterfaceType reflect.Type

func init() {
	hasIdInterfaceType = reflect.TypeOf((*HasId)(nil)).Elem()

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

func (r *Repo) SaveNew(item interface{}) (bson.ObjectId, error) {
	session := r.session.Copy()
	defer session.Close()

	coll := session.DB("").C(r.coll)

	wrapper := wrap(item, bson.NewObjectId())
	err := coll.Insert(wrapper)

	return wrapper.Id, err
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
	wrapper := wrapper{}
	err := coll.Find(bson.M{"_id": id}).One(&wrapper)
	if err != nil {
		return nil, err
	}

	i, err := unwrap(wrapper)
	if err != nil {
		return nil, errors.Wrap(err, "Unwrapping generated an error")
	}

	return i, nil
}

func (r *Repo) GetAll(result interface{}) error {

	resultv := reflect.ValueOf(result)
	if resultv.Kind() != reflect.Ptr {
		panic("result argument must be a slice address")
	}

	slicev := resultv.Elem()

	if slicev.Kind() == reflect.Interface {
		slicev = slicev.Elem()
	}
	if slicev.Kind() != reflect.Slice {
		panic("result argument must be a slice address")
	}

	session := r.session.Copy()
	defer session.Close()

	coll := session.DB("").C(r.coll)

	iter := coll.Find(bson.M{}).Iter()

	var errs = make([]string, 0)
	var w wrapper
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

func wrap(item interface{}, id bson.ObjectId) wrapper {

	//util.Introspect(item)

	t, v := util.ConcreteType(item)

	log.Println()
	log.Println(">>>", t, v)
	log.Println()

	//setId(item, id)

	// The value is saved as a wrapper value, with V being the actual value and T being its type key.
	wrapper := wrapper{
		Id: id,
		T:  typeregistry.GetKey(t),
		V:  item,
	}
	return wrapper
}

func setId(item interface{}, id bson.ObjectId) {

	v := util.ConcreteValue(item)
	log.Print(v)

	util.Introspect(v)

	pv := v.Addr()

	if pv.Type().Implements(hasIdInterfaceType) {

		hasId := pv.Interface().(HasId)
		hasId.SetId(id.Hex())
		log.Println()
		log.Println(">>>>>>", hasId)
		log.Println()
	}

}

func unwrap(w wrapper) (interface{}, error) {
	//tv := reflect.TypeOf(w.V)
	//fmt.Println(tv)

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
