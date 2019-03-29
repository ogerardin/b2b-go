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

func ConcreteValue(v interface{}) reflect.Value {
	log.Tracef("%T %[1]V", v)
	rv := reflect.ValueOf(v)
	log.Tracef("%T %[1]V", rv)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
		log.Tracef("%T %[1]V", rv)
	}
	return rv
}

func print(context string, rv reflect.Value) {
	fmt.Println(">>>"+context, rv.Kind(), rv.Type())
}
