package typeregistry

import "reflect"

var types map[string]reflect.Type

func init() {
	types = make(map[string]reflect.Type)
}

func Register(t reflect.Type) {
	key := GetKey(t)
	types[key] = t
}

func GetKey(t reflect.Type) string {
	key := t.PkgPath() + "." + t.Name()
	return key
}

func GetType(key string) reflect.Type {
	t := types[key]
	return t
}
