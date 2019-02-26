package domain

import (
	"b2b-go/lib/typeregistry"
	"fmt"
	"reflect"
)

//Represents a remote peer backup destination.
type PeerTarget struct {
	BackupTargetBase
}

func (t PeerTarget) Desc() string {
	return fmt.Sprintf("Peer Target (%+v)", t)
}

func init() {
	typeregistry.Register(reflect.TypeOf((*PeerTarget)(nil)).Elem())
}