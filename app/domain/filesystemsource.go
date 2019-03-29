package domain

import (
	"b2b-go/lib/typeregistry"
	"fmt"
	"reflect"
)

type FilesystemSource struct {
	BackupSourceBase `mapstructure:",squash"`
	Paths            []string
}

var _ BackupSource = &FilesystemSource{}

func (fss FilesystemSource) Desc() string {
	return fmt.Sprintf("File system source (%+v)", fss)
}

func init() {
	typeregistry.Register(reflect.TypeOf((*FilesystemSource)(nil)).Elem())
}
