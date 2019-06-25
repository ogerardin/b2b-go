package domain

import (
	"fmt"
)

type BackupTarget interface {
	backupTarget()
	Desc() string
}

type BackupTargetBase struct {
	Id      string `mapstructure:"_id" json:"id"`
	Enabled bool   `json:"enabled"`
	Name    string `json:"name"`
}

func (BackupTargetBase) backupTarget() {}

func (btb BackupTargetBase) Desc() string {
	return fmt.Sprintf("BackupTargetBase (%+v)", btb)
}
