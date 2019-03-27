package domain

import (
	"fmt"
)

type BackupSource interface {
	Desc() string
	//SetId(id string)
}

type BackupSourceBase struct {
	Id      string
	Enabled bool
	Name    string
}

var _ BackupSource = &BackupSourceBase{}

func (bsb *BackupSourceBase) GetId() string {
	return bsb.Id
}

func (bsb *BackupSourceBase) SetId(id string) {
	bsb.Id = id
}

func (bsb BackupSourceBase) Desc() string {
	return fmt.Sprintf("BackupSourceBase (%+v)", bsb)
}
