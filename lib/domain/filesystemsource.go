package domain

type FilesystemSource struct {
	BackupSource
	Paths []string
}
