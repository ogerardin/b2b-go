package domain

import (
	"b2b-go/lib/typeregistry"
	"fmt"
	"reflect"
)

//Represents the local (internal) backup destination.
type LocalTarget struct {
	BackupTargetBase
}

var _ BackupTarget = &LocalTarget{}

func (t LocalTarget) Desc() string {
	return fmt.Sprintf("Local Target (%+v)", t)
}

func init() {
	typeregistry.Register(reflect.TypeOf((*LocalTarget)(nil)).Elem())
}
