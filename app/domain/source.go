package domain

import "fmt"

type BackupSource interface {
	backupSource()
	Desc() string
}

type BackupSourceBase struct {
	Enabled bool
	Name    string
}

func (bsb BackupSourceBase) Desc() string {
	return fmt.Sprintf("BackupSourceBase (%+v)", bsb)
}

func (BackupSourceBase) backupSource() {}
