package util

import (
	"reflect"
)

func ConcreteValue(v interface{}) reflect.Value {
	rv := reflect.ValueOf(v)
	//	log.Print(rv.Type(), rv.Kind(), rv)

	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
		//		log.Print(rv.Type(), rv.Kind(), rv)
	}
	return rv
}
