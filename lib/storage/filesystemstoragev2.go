package storage

import (
	"b2b-go/lib/util"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"path/filepath"
)

var _ StorageService = (*FilesystemStorageV2)(nil)

type FilesystemStorageV2 struct {
	FilesystemStorage
}

func (s *FilesystemStorageV2) StoreReader(f io.Reader, filename string) (string, error) {
	localPath := s.getLocalPathForNextRevision(filename)
	log.Printf("Writing revision to %v", localPath)

	dir := filepath.Dir(localPath)
	err := os.MkdirAll(dir, os.ModeDir|util.OS_USER_RWX|util.OS_ALL_R)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to create directory %s", dir)
	}

	target, err := os.OpenFile(localPath, os.O_CREATE|os.O_WRONLY, 0)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to create local file %s", localPath)
	}
	_, err = io.Copy(target, f)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to copy bytes to local file")
	}
	revisionId, _ := filepath.Rel(s.baseDirectory, localPath)
	return revisionId, nil
}

func (s *FilesystemStorageV2) getLocalPathForNextRevision(filename string) string {
	localfile := s.remoteToLocal(filename)
	base := filepath.Base(localfile)
	dir := filepath.Dir(localfile)

	for revNum := 0; ; revNum++ {
		revisionFilename := buildRevisionFilename(base, revNum)
		revisionFilePath := filepath.Join(dir, revisionFilename)
		if exists, _ := util.FileExists(revisionFilePath); exists {
			return revisionFilePath
		}
	}
}

func (s *FilesystemStorageV2) GetAllFiles() ([]FileInfo, error) {
	revisions := s.GetAllRevisions()

	latestRevisionByFile := make(map[string]RevisionInfo, 0)
	for _, rev := range revisions {
		filename := rev.filename
		currentrev, found := latestRevisionByFile[filename]
		if !found || rev.storedDate.After(currentrev.storedDate) {
			latestRevisionByFile[filename] = rev
		}
	}

	result := make([]FileInfo, 0)
	for _, rev := range latestRevisionByFile {
		result = append(result, *rev.getFileInfo())
	}

	return result, nil
}

func (s *FilesystemStorageV2) GetAllRevisions() []RevisionInfo {
	result := make([]RevisionInfo, 0)
	err := filepath.Walk(s.baseDirectory, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() {
			info, err := s.buildRevisionInfo(path)
			if err != nil {
				return err
			}
			result = append(result, *info)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return result
}

func (s *FilesystemStorageV2) buildRevisionInfo(localpath string) (*RevisionInfo, error) {
	info, err := s.FilesystemStorage.buildRevisionInfo(localpath)
	if err != nil {
		return nil, err
	}

	relpath, err := filepath.Rel(s.baseDirectory, localpath)
	if err != nil {
		return nil, err
	}

	info.id = relpath
	return info, nil
}

func buildRevisionFilename(base string, revNum int) string {
	return fmt.Sprintf("%s#%d", base, revNum)
}

func NewV2(baseDirectory string) (*FilesystemStorageV2, error) {
	v1, _ := New(baseDirectory)
	return &FilesystemStorageV2{
		*v1,
	}, nil
}
