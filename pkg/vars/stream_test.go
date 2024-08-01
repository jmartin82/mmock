package vars

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Go\n"))
	}))
	defer server.Close()

	k := fmt.Sprintf("http.contents(%s)", server.URL)
	holders := []string{k}

	result := st.Fill(holders)
	v, f := result[k]
	if !f {
		t.Errorf("Stream key not found")
	}

	if strings.TrimSpace(v[0]) != "Go" {
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
