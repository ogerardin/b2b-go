package storage

import (
	"io"
)

type StorageService interface {
	Store(file string) (string, error)
	GetAsReader(key string) (io.Reader, error)
	Delete(key string) error
	GetAllFiles() ([]string, error)
}
