package domain

import "fmt"

type BackupTarget interface {
	backupTarget()
	Desc() string
}

type BackupTargetBase struct {
	Enabled bool
	Name    string
}

func (BackupTargetBase) backupTarget() {}

func (btb *BackupTargetBase) SetId(id string) {
	btb.SetId(id)
}

func (btb BackupTargetBase) Desc() string {
	return fmt.Sprintf("BackupTargetBase (%+v)", btb)
}
