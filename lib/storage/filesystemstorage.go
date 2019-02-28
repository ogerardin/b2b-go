package storage

import (
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
)

var _ StorageService = (*FilesystemStorage)(nil)

type FilesystemStorage struct {
	baseDirectory string
}

func (s *FilesystemStorage) GetAllRevisions() []RevisionInfo {
	panic("implement me")
}

func (s *FilesystemStorage) localToRemote(localpath string) string {
	relativePath, _ := filepath.Rel(s.baseDirectory, localpath)
	const sep = string(filepath.Separator)
	return filepath.Join(sep, relativePath)

}

func (s *FilesystemStorage) GetAllFiles() ([]string, error) {
	result := make([]string, 0)
	err := filepath.Walk(s.baseDirectory, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			result = append(result, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *FilesystemStorage) Store(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	return s.StoreReader(f, filename)
}

func (s *FilesystemStorage) StoreReader(f io.Reader, filename string) (string, error) {
	localPath := s.remoteToLocal(filename)

	dir := filepath.Dir(localPath)
	err := os.MkdirAll(dir, os.ModeDir)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to create directory %s", dir)
	}

	target, err := os.OpenFile(localPath, os.O_CREATE|os.O_WRONLY, 0)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(target, f)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (s *FilesystemStorage) remoteToLocal(filename string) string {
	filename, _ = filepath.Abs(filename)
	volumeName := filepath.VolumeName(filename)
	filename = filename[len(volumeName):]
	return filepath.Join(s.baseDirectory, filename)
}

func (s *FilesystemStorage) GetAsReader(filename string) (io.Reader, error) {
	localPath := s.remoteToLocal(filename)
	return s.getAsReader(localPath, filename)
}

func (*FilesystemStorage) Delete(filename string) error {
	return os.Remove(filename)
}

func (s *FilesystemStorage) getAsReader(localPath string, filename string) (io.Reader, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func New(baseDirectory string) (*FilesystemStorage, error) {
	fileInfo, err := os.Stat(baseDirectory)
	if !fileInfo.IsDir() {
		return nil, NewStorageError(baseDirectory + " is not a directory")
	}
	if err != nil {
		return nil, err
	}
	return &FilesystemStorage{
		baseDirectory: baseDirectory,
	}, nil
}
