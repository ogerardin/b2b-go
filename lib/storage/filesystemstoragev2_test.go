package storage

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFilesystemStoragev2(t *testing.T) {
	d, _ := ioutil.TempDir(os.TempDir(), "mongotools-test")
	fss, _ := NewV2(d)
	StoreAndretrieveMultipleRevisionsTest(t, fss)
}
