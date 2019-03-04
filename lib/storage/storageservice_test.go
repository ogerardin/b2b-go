package storage

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func getSampleFilePaths() []string {
	result := make([]string, 0)
	err := filepath.Walk("../../test-fileset", func(path string, f os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if !f.IsDir() {
			result = append(result, path)
		}
		return nil
	})
	if err != nil {
		panic("Failed to list sample test files")
	}
	return result
}

func StoreAndRetrieveTest(t *testing.T, s StorageService) {
	paths0 := getSampleFilePaths()

	for _, p := range paths0 {
		_, err := s.Store(p)
		if err != nil {
			t.Fatal(err)
		}
	}

	paths1, err := s.GetAllFiles()
	if err != nil {
		t.Fatal(err)
	}

	sort.Strings(paths0)
	sort.Strings(paths1)

	for i, p0 := range paths0 {
		p1 := paths1[i]
		assertMatches(t, s, p0, p1)
	}
}

func StoreAndretrieveMultipleRevisionsTest(t *testing.T, s StorageService) {
	paths0 := getSampleFilePaths()

	tempFile, _ := ioutil.TempFile("", "*.bin")
	tempFile.Close()
	tempFileName := tempFile.Name()

	for _, f := range paths0 {
		t.Logf("Creating new revision of %s", f)
		copyFile(f, tempFileName)
		_, err := s.Store(tempFileName)
		if err != nil {
			t.Fatal(err)
		}
	}

	revisions := s.GetAllRevisions()
	sort.Slice(revisions, func(i, j int) bool {
		ri := revisions[i]
		rj := revisions[j]
		return ri.storedDate.Before(rj.storedDate)
	})
	t.Logf("Retrieved revision list: %v", revisions)

	assert.Equal(t, len(paths0), len(revisions))

	for i, f := range paths0 {
		rev := revisions[i]
		revisionId := rev.id
		assertStoredRevisionMatchesFile(t, s, f, revisionId)

	}
}

func assertStoredRevisionMatchesFile(t *testing.T, s StorageService, f0 string, id string) {
	bytes0, err := ioutil.ReadFile(f0)
	if err != nil {
		t.Fatal(err)
	}
	reader1, err := s.GetAsReader(id)
	if err != nil {
		t.Fatal(err)
	}
	bytes1, err := ioutil.ReadAll(reader1)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(bytes0, bytes1) != 0 {
		t.Fatal("files are different")
	}

}

func copyFile(s string, t string) {
	fs, err := os.Open(s)
	if err != nil {
		panic(err)
	}
	ft, err := os.Create(t)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(ft, fs)
	if err != nil {
		panic(err)
	}

}

func assertMatches(t *testing.T, s StorageService, f0 string, f1 string) {
	bytes0, err := ioutil.ReadFile(f0)
	if err != nil {
		t.Fatal(err)
	}
	reader1, err := s.GetAsReader(f1)
	if err != nil {
		t.Fatal(err)
	}
	bytes1, err := ioutil.ReadAll(reader1)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(bytes0, bytes1) != 0 {
		t.Fatal("files are different")
	}
}
