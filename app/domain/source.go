package domain

import (
	"fmt"
)

type BackupSource interface {
	backupSource()
	Desc() string
}

type BackupSourceBase struct {
	Id      string `mapstructure:"_id" json:"-" jsonapi:"primary,sources"`
	Enabled bool   `json:"enabled"`
	Name    string `json:"name"`
}

func (bsb BackupSourceBase) GetID() string {
	return bsb.Id
}

func (bsb BackupSourceBase) backupSource() {}

func (bsb BackupSourceBase) Desc() string {
	return fmt.Sprintf("BackupSourceBase (%+v)", bsb)
}
