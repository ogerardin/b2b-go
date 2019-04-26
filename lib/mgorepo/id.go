package mgorepo

import "github.com/globalsign/mgo/bson"

type IdGenerator interface {
	NewId() interface{}
}

type ObjectIdGenerator struct{}

func (*ObjectIdGenerator) NewId() interface{} {
	return bson.NewObjectId().Hex()
}

func defaultGenerator() IdGenerator {
	return &ObjectIdGenerator{}
}
