package domain

import (
	"b2b-go/app"
	"b2b-go/lib/typeregistry"
	"fmt"
	"reflect"
)

//Represents a remote peer backup destination.
type PeerTarget struct {
	BackupTargetBase
	Hostname string
	Port     int
}

func (t PeerTarget) Desc() string {
	return fmt.Sprintf("Peer Target (%+v)", t)
}

func init() {
	typeregistry.Register(reflect.TypeOf((*PeerTarget)(nil)).Elem())
}