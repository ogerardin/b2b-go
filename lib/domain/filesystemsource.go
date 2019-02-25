package domain

import (
	"b2b-go/lib/typeregistry"
	"fmt"
	"reflect"
)

type FilesystemSource struct {
	BackupSourceBase
	Paths []string
}

func (fss FilesystemSource) Desc() string {
	return fmt.Sprintf("File system source (%+v)", fss)
}

func init() {
	typeregistry.Register(reflect.TypeOf((*FilesystemSource)(nil)).Elem())
}
