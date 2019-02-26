package storage

import (
	"testing"
)

func TestFilesystemStorage(t *testing.T) {
	fss, _ := New(".")
	StoreAndRetrieve(t, fss)
}
