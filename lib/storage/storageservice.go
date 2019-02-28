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
}

type RevisionInfo struct {
	id         string
	filename   string
	storedDate time.Time
	size       uint
	deleted    bool
}
