package typeregistry

import (
	"log"
	"reflect"
)

var types map[string]reflect.Type

func init() {
	types = make(map[string]reflect.Type)
}

func Register(t reflect.Type) {
	key := GetKey(t)
	types[key] = t
}

func GetKey(t reflect.Type) string {
	pkgPath := t.PkgPath()
	name := t.Name()
	if len(pkgPath) == 0 || len(name) == 0 {
		log.Panic("non-defined type")
	}

	key := pkgPath + "." + name
	return key
}

func GetType(key string) reflect.Type {
	t := types[key]
	return t
}
