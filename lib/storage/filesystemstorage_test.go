package storage

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFilesystemStorage(t *testing.T) {
	d, _ := ioutil.TempDir(os.TempDir(), "mongotools-test")
	fss, _ := New(d)
	StoreAndRetrieve(t, fss)
}
