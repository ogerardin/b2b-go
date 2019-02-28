package storage

import (
	"bytes"
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

func StoreAndRetrieve(t *testing.T, s StorageService) {
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
