package mgorepo

import (
	"b2b-go/lib/typeregistry"
	"b2b-go/lib/util"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"reflect"
	"runtime"
	"strings"
)

//FIXME loggers should be configured centrally
var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			funcVal := frame.Function
			return funcVal, ""
		},
	})
	log.SetLevel(logrus.TraceLevel)
}

type Repo struct {
	session *mgo.Session
	coll    string
}

type HasGetId interface {
	GetId() string
}

type HasSetId interface {
	SetId(id string)
}

// this should be a const but compiler complains "const initializer is not a constant"
var hasIdInterfaceType reflect.Type

func init() {
	hasIdInterfaceType = reflect.TypeOf((*HasSetId)(nil)).Elem()
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
	t := util.ConcreteValue(item).Type()

	setId(item, id)

	// The value is saved as a wrapper value, with V being the actual value and T being its type key.
	wrapper := wrapper{
		Id: id,
		T:  typeregistry.GetKey(t),
		V:  item,
	}
	return wrapper
}

func setId(item interface{}, id bson.ObjectId) {

	defer util.RecoverPanicAndLog(log, "setId failed")

	log.Tracef("item = %T %[1]+v", item)

	v := reflect.ValueOf(item)
	log.Tracef("value = %v %v %+v", v.Kind(), v.Type(), v)

	// if we don't have a pointer, we won't be able to call SetId
	if !(v.Kind() == reflect.Ptr) {
		log.Debug("not a pointer")
		return
	}

	// get concrete structure
	for v.Kind() != reflect.Struct {
		v = v.Elem()
		log.Tracef("elem = %v %v %+v", v.Kind(), v.Type(), v)
	}

	// obtain pointer to this structure
	if !v.CanAddr() {
		log.Debug("not an addressable value")
		return
	}
	pv := v.Addr()
	log.Tracef("ptr = %v %v %+v", pv.Kind(), pv.Type(), pv)

	// if it does not implement HasSetId, nothing to do
	if !pv.Type().Implements(hasIdInterfaceType) {
		return
	}

	// convert to HasSetId type and call SetId
	itemWithId := pv.Interface().(HasSetId)
	log.Tracef("HasSetId = %T %[1]V", itemWithId)
	itemWithId.SetId(id.Hex())

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
	return i, nil

}
