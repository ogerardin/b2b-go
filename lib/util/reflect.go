package util

import (
	"fmt"
	"log"
	"reflect"
)

func Introspect(v interface{}) {
	rv := reflect.ValueOf(v)
	context := ""
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		print(context, rv)
		rv = rv.Elem()
		context = context + "  "
	}
	print(context, rv)
}

func print(context string, rv reflect.Value) {
	fmt.Println(">>>"+context, rv.Kind(), rv.Type())
}

func ConcreteValue(v interface{}) reflect.Value {
	rv := reflect.ValueOf(v)
	log.Print(rv.Type(), rv.Kind(), rv)

	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
		log.Print(rv.Type(), rv.Kind(), rv)
	}
	return rv
}
