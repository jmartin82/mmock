package definition

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var updatesCh = make(chan []Mock)

type mockReader struct {
	canRead int
	readOk  int
}

func (mr *mockReader) CanRead(filename string) bool {
	if mr.canRead > 0 {
		mr.canRead--
		return true
	}
	return false
}

func (mr *mockReader) Read(content []byte) (Mock, error) {
	if mr.readOk > 0 {
		mr.readOk--

		m := Mock{}
		m.Control.Priority = mr.readOk
		return m, nil
	}

	return Mock{}, errors.New("error")

}

func TestReadMocksDefinition(t *testing.T) {

	content := []byte("temporary file's content")
	dir, err := ioutil.TempDir("", "mmock")
	if err != nil {
		t.Errorf("Error creating temporary folder")
	}

	numDefFiles := 5
	for i := 0; i < numDefFiles; i++ {
		tmpfn := filepath.Join(dir, fmt.Sprintf("tmpfile_%d", i))
		if err := ioutil.WriteFile(tmpfn, content, 0666); err != nil {
			t.Errorf("Error creating temporary file")
		}
	}

	defer os.RemoveAll(dir) // clean up

	fileDef := NewFileDefinition(dir, updatesCh)
	fileDef.AddConfigReader(&mockReader{canRead: 3, readOk: 2})
	mocks := fileDef.ReadMocksDefinition()
	if len(mocks) != 2 {
		t.Errorf("Error getting mocks definition, expected %d got %d", numDefFiles, len(mocks))
	}
}

func TestPriority(t *testing.T) {

	content := []byte("temporary file's content")
	dir, err := ioutil.TempDir("", "mmock1")
	if err != nil {
		t.Errorf("Error creating temporary folder")
	}

	numDefFiles := 5
	for i := 0; i < numDefFiles; i++ {
		tmpfn := filepath.Join(dir, fmt.Sprintf("tmpfile_%d", i))
		if err := ioutil.WriteFile(tmpfn, content, 0666); err != nil {
			t.Errorf("Error creating temporary file")
		}
	}

	defer os.RemoveAll(dir) // clean up

	fileDef := NewFileDefinition(dir, updatesCh)
	fileDef.AddConfigReader(&mockReader{canRead: 5, readOk: 5})
	mocks := fileDef.ReadMocksDefinition()
	for i, m := range mocks {
		if (4 - i) != m.Control.Priority {
			t.Errorf("Expected priority %d got %d", (5 - i), m.Control.Priority)
		}
	}

}

func TestHotReplace(t *testing.T) {
	content := []byte("temporary file's content")
	dir, err := ioutil.TempDir("", "mmock2")
	if err != nil {
		t.Errorf("Error creating temporary folder")
	}

	tmpfn := filepath.Join(dir, "tmpfile")
	if err := ioutil.WriteFile(tmpfn, content, 0666); err != nil {
		t.Errorf("Error creating temporary file")
	}

	defer os.RemoveAll(dir) // clean up

	fileDef := NewFileDefinition(dir, updatesCh)
	fileDef.AddConfigReader(&mockReader{canRead: 5, readOk: 5})
	fileDef.WatchDir()

	if err := ioutil.WriteFile(tmpfn, content, 0666); err != nil {
		t.Errorf("Error updating temporary file")
	}

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(5 * time.Second)
		timeout <- true
	}()
	select {
	case <-updatesCh:
		t.Logf("New channel definition")
	case <-timeout:
		t.Fail()
	}

	fileDef.UnWatchDir()
	close(updatesCh)

}

/*
func TestUpdateOnFileChange(t *testing.T) {
	const s, sep, want = "chicken", "ken", 4
	t.Errorf("Index(%q,%q) = %v; want %v", s, sep, got, want)

}
*/
