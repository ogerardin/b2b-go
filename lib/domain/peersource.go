package domain

import (
	"b2b-go/lib/typeregistry"
	"fmt"
	"github.com/google/uuid"
	"reflect"
)

type PeerSource struct {
	BackupSourceBase
	remoteCOmputerId uuid.UUID
}

func (ps PeerSource) Desc() string {
	return fmt.Sprintf("Peer source (%+v)", ps)
}

func init() {
	typeregistry.Register(reflect.TypeOf((*PeerSource)(nil)).Elem())
}
