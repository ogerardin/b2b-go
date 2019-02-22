package domain

type BackupSource interface {
	backupSource()
}

type BackupSourceBase struct {
	Enabled bool
	Name    string
}

func (BackupSourceBase) backupSource() {}
