package util

import (
	"fmt"
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

func ConcreteType(v interface{}) reflect.Type {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}
	return rv.Type()
}

func print(context string, rv reflect.Value) {
	fmt.Println(">>>"+context, rv.Kind(), rv.Type())
}