package domain

import (
	"b2b-go/lib/typeregistry"
	"fmt"
	"reflect"
)

//Represents a local filesystem backup destination
type FilesystemTarget struct {
	BackupTargetBase `mapstructure:",squash"`
	Path             string
}

var _ BackupTarget = &FilesystemTarget{}

func (t FilesystemTarget) Desc() string {
	return fmt.Sprintf("Filesystem Target (%+v)", t)
}

func init() {
	typeregistry.Register(reflect.TypeOf((*FilesystemTarget)(nil)).Elem())
}
