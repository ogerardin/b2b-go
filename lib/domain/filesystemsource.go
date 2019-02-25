package domain

import (
	tr "b2b-go/lib/typeregistry"
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
	tr.Register(reflect.TypeOf((*FilesystemSource)(nil)).Elem())
}
