package storage

import (
	"io"
	"time"
)

type StorageService interface {
	Store(file string) (string, error)
	StoreReader(f io.Reader, filename string) (string, error)
	GetAsReader(key string) (io.Reader, error)
	Delete(key string) error
	GetAllFiles() ([]string, error)
	GetAllRevisions() []RevisionInfo
	localToRemote(path string) string
	buildRevisionInfo(localpath string) (*RevisionInfo, error)
}

type RevisionInfo struct {
	id         string
	filename   string
	storedDate time.Time
	size       int64
	deleted    bool
}

type FileInfo struct {
	path    string
	deleted bool
}

func (ri *RevisionInfo) getFileInfo() *FileInfo {
	return &FileInfo{
		path:    ri.filename,
		deleted: ri.deleted,
	}
}
