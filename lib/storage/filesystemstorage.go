package storage

import (
	"b2b-go/lib/util"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
)

var _ StorageService = (*FilesystemStorage)(nil)

type FilesystemStorage struct {
	baseDirectory string
}

func (s *FilesystemStorage) GetRevisionAsReader(key string) (io.Reader, error) {
	panic("not implemented")
}

func (s *FilesystemStorage) buildRevisionInfo(localpath string) (*RevisionInfo, error) {
	info, err := os.Stat(localpath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed top get file info for %s", localpath)
	}
	revisionInfo := RevisionInfo{
		id:         "",
		filename:   localpath,
		size:       info.Size(),
		storedDate: info.ModTime(),
	}
	return &revisionInfo, nil
}

func (s *FilesystemStorage) GetAllRevisions() []RevisionInfo {
	panic("implement me")
}

func (s *FilesystemStorage) localToRemote(localpath string) string {
	relativePath, _ := filepath.Rel(s.baseDirectory, localpath)
	const sep = string(filepath.Separator)
	return filepath.Join(sep, relativePath)

}

func (s *FilesystemStorage) GetAllFiles() ([]FileInfo, error) {
	result := make([]FileInfo, 0)
	err := filepath.Walk(s.baseDirectory, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fileInfo := FileInfo{
				Path:    path,
				Deleted: false,
			}
			result = append(result, fileInfo)
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
	err := os.MkdirAll(dir, os.ModeDir|util.OS_USER_RWX|util.OS_ALL_R)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to create directory %s", dir)
	}

	target, err := os.OpenFile(localPath, os.O_CREATE|os.O_WRONLY, util.OS_USER_RWX|util.OS_ALL_R)
	if err != nil {
		return "", errors.Wrapf(err, "failed to open target file for writing: %s", localPath)
	}
	_, err = io.Copy(target, f)
	if err != nil {
		return "", errors.Wrapf(err, "failed to copy to: %s", localPath)
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
