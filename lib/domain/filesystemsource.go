package domain

import (
	tr "b2b-go/lib/typeregistry"
	"reflect"
)

type FilesystemSource struct {
	BackupSourceBase
	Paths []string
}

func init() {
	tr.Register(reflect.TypeOf((*FilesystemSource)(nil)).Elem())
}
