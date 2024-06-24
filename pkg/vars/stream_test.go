package vars

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"path/filepath"
	"strings"
	"testing"
	"github.com/jmartin82/mmock/v3/internal/dummyServer"
)

func TestReadFile(t *testing.T) {

	content := []byte("This is a big file")
	dir, err := ioutil.TempDir("", "mmock")
	if err != nil {
		t.Errorf("Error creating temporary folder")
	}

	tmpfn := filepath.Join(dir, "bigfile")
	if err := ioutil.WriteFile(tmpfn, content, 0666); err != nil {
		t.Errorf("Error updating temporary file")
	}

	defer os.RemoveAll(dir) // clean up

	st := Stream{}

	k := fmt.Sprintf("file.contents(%s)", tmpfn)
	holders := []string{k}

	result := st.Fill(holders)
	v, f := result[k]
	if !f {
		t.Errorf("Stream key not found")
	}

	if !strings.Contains(v[0], "This is a big file") {
		t.Errorf("Couldn't get the content. Value: %s", v)
	}

}

func TestHTTPContent(t *testing.T) {
	st := Stream{}
	wg := &sync.WaitGroup{}
	wg.Add(1)

	dums := dummyServer.Start(wg, 8937)
	k := "http.contents(http://localhost:8937/hello)"
	holders := []string{k}

	result := st.Fill(holders)
	dums.Stop()
	v, f := result[k]
	if !f {
		t.Errorf("Stream key not found")
	}

	if !strings.Contains(v[0], "hello") {
		t.Errorf("Couldn't get the content. Value: %s", v)
	}
}

func TestError(t *testing.T) {
	st := Stream{}

	k := "file.contents(XXXXX)"
	holders := []string{k}

	result := st.Fill(holders)
	v, f := result[k]
	if !f {
		t.Errorf("Stream key not found")
	}

	if !strings.Contains(v[0], "ERROR: open XXXXX: no such file or directory") {
		t.Errorf("Couldn't get the content. Value: %s", v)
	}
}
