package definition

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestHotReplace(t *testing.T) {

	content := []byte("temporary file's content")
	dir, err := ioutil.TempDir("", "mmock2")
	if err != nil {
		t.Errorf("Error creating temporary folder")
	}

	tmpfn := filepath.Join(dir, "tmpfile")

	defer os.RemoveAll(dir) // clean up

	fsUpdate := make(chan struct{})
	mapper := NewFileWatcher(dir, fsUpdate)

	mapper.Bind()

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(5 * time.Second)
		timeout <- true
	}()

	go func() {
		time.Sleep(1 * time.Second)
		if err := ioutil.WriteFile(tmpfn, content, 0666); err != nil {
			t.Errorf("Error updating temporary file")
		}
	}()

	select {
	case <-fsUpdate:
		t.Logf("New channel definition")
	case <-timeout:
		t.Fail()
	}

	mapper.UnBind()
	close(fsUpdate)

}
