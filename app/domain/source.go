package domain

import (
	"fmt"
)

type BackupSource interface {
	backupSource()
	Desc() string
}

type BackupSourceBase struct {
	Id      string `mapstructure:"_id" json:"id"`
	Enabled bool   `json:"enabled"`
	Name    string `json:"name"`
}

func (bsb BackupSourceBase) backupSource() {}

func (bsb BackupSourceBase) Desc() string {
	return fmt.Sprintf("BackupSourceBase (%+v)", bsb)
}
