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

func (s *FilesystemStorageV2) GetAllFiles() ([]string, error) {
	_ = s.GetAllRevisions()
	//TODO
	return nil, nil
}

func (s *FilesystemStorageV2) GetAllRevisions() []RevisionInfo {
	result := make([]RevisionInfo, 0)
	err := filepath.Walk(s.baseDirectory, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if !f.IsDir() {
			result = append(result, s.buildRevisionInfo(path))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return result
}

func (s *FilesystemStorageV2) buildRevisionInfo(localpath string) RevisionInfo {
	//TODO
	return RevisionInfo{}
}

func buildRevisionFilename(base string, revNum int) string {
	return fmt.Sprintf("%s#%d", base, revNum)
}
