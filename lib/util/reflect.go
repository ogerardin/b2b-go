package util

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	log.SetLevel(logrus.TraceLevel)
}

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
	log.Trace(rv.Type(), rv.Kind(), rv)

	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
		log.Trace(rv.Type(), rv.Kind(), rv)
	}
	return rv
}
