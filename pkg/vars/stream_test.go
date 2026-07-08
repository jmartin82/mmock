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

	"github.com/jmartin82/mmock/v3/pkg/mock"
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

func TestStreamPathParamOnly(t *testing.T) {
	content := []byte("csv content for arquivo")
	dir, err := ioutil.TempDir("", "mmock")
	if err != nil {
		t.Errorf("Error creating temporary folder")
	}
	defer os.RemoveAll(dir) // clean up

	if err := ioutil.WriteFile(filepath.Join(dir, "arquivo.csv"), content, 0666); err != nil {
		t.Errorf("Error creating temporary file")
	}

	req := &mock.Request{Path: "/results/arquivo.csv"}
	def := &mock.Definition{Request: mock.Request{Path: "/results/:filename"}}
	st := Stream{Mock: def, Request: req}

	k := fmt.Sprintf("file.contents(%s/{{request.path.filename}})", dir)
	result := st.Fill([]string{k})

	if !strings.Contains(result[k][0], "csv content for arquivo") {
		t.Errorf("Path param was not resolved. Value: %s", result[k])
	}
}

func TestStreamBodyJSONOnly(t *testing.T) {
	content := []byte("body json content")
	dir, err := ioutil.TempDir("", "mmock")
	if err != nil {
		t.Errorf("Error creating temporary folder")
	}
	defer os.RemoveAll(dir) // clean up

	if err := ioutil.WriteFile(filepath.Join(dir, "exemplo.csv"), content, 0666); err != nil {
		t.Errorf("Error creating temporary file")
	}

	req := &mock.Request{Path: "/results"}
	req.Body = `{"name":"exemplo"}`
	req.Headers = mock.Values{"Content-Type": {"application/json"}}
	def := &mock.Definition{Request: mock.Request{Path: "/results"}}
	st := Stream{Mock: def, Request: req}

	k := fmt.Sprintf("file.contents(%s/{{request.body.name}}.csv)", dir)
	result := st.Fill([]string{k})

	if !strings.Contains(result[k][0], "body json content") {
		t.Errorf("Body JSON value was not resolved. Value: %s", result[k])
	}
}

func TestStreamCombinedPathAndBody(t *testing.T) {
	content := []byte("combined content")
	dir, err := ioutil.TempDir("", "mmock")
	if err != nil {
		t.Errorf("Error creating temporary folder")
	}
	defer os.RemoveAll(dir) // clean up

	if err := os.MkdirAll(filepath.Join(dir, "acme", "exemplos"), 0777); err != nil {
		t.Errorf("Error creating nested folder")
	}
	if err := ioutil.WriteFile(filepath.Join(dir, "acme", "exemplos", "exemplo.csv"), content, 0666); err != nil {
		t.Errorf("Error creating temporary file")
	}

	req := &mock.Request{Path: "/blu365-athena-results/acme/athena-results/x.csv"}
	req.Body = `{"name":"exemplo"}`
	req.Headers = mock.Values{"Content-Type": {"application/json"}}
	def := &mock.Definition{Request: mock.Request{Path: "/blu365-athena-results/:tenant/athena-results/:filename"}}
	st := Stream{Mock: def, Request: req}

	k := fmt.Sprintf("file.contents(%s/{{request.path.tenant}}/exemplos/{{request.body.name}}.csv)", dir)
	result := st.Fill([]string{k})

	if !strings.Contains(result[k][0], "combined content") {
		t.Errorf("Combined path+body was not resolved. Value: %s", result[k])
	}
}

func TestStreamQueryParam(t *testing.T) {
	content := []byte("query content")
	dir, err := ioutil.TempDir("", "mmock")
	if err != nil {
		t.Errorf("Error creating temporary folder")
	}
	defer os.RemoveAll(dir) // clean up

	if err := ioutil.WriteFile(filepath.Join(dir, "exemplo.csv"), content, 0666); err != nil {
		t.Errorf("Error creating temporary file")
	}

	req := &mock.Request{Path: "/results"}
	req.QueryStringParameters = mock.Values{"name": {"exemplo"}}
	def := &mock.Definition{Request: mock.Request{Path: "/results"}}
	st := Stream{Mock: def, Request: req}

	k := fmt.Sprintf("file.contents(%s/{{request.query.name}}.csv)", dir)
	result := st.Fill([]string{k})

	if !strings.Contains(result[k][0], "query content") {
		t.Errorf("Query param was not resolved. Value: %s", result[k])
	}
}

func TestStreamBodyValueTraversalRejected(t *testing.T) {
	// A body value resolving to a traversal token must not be substituted; the literal
	// holder is kept so no parent directory can be reached.
	req := &mock.Request{Path: "/results"}
	req.Body = `{"name":".."}`
	req.Headers = mock.Values{"Content-Type": {"application/json"}}
	def := &mock.Definition{Request: mock.Request{Path: "/results"}}
	st := Stream{Mock: def, Request: req}

	k := "file.contents(/data/base/{{request.body.name}}.csv)"
	result := st.Fill([]string{k})

	if !strings.Contains(result[k][0], "{{request.body.name}}") {
		t.Errorf("Traversal value was substituted instead of rejected. Value: %s", result[k])
	}
}

func TestStreamGjsonRegexArgParens(t *testing.T) {
	// The argument contains gjson .regex(...)/.concat(...) with parentheses and a {4}
	// quantifier; it must resolve via the gjson chain (README example semantics).
	content := []byte("regex chain content")
	dir, err := ioutil.TempDir("", "mmock")
	if err != nil {
		t.Errorf("Error creating temporary folder")
	}
	defer os.RemoveAll(dir) // clean up

	if err := ioutil.WriteFile(filepath.Join(dir, "2307-x.csv"), content, 0666); err != nil {
		t.Errorf("Error creating temporary file")
	}

	req := &mock.Request{Path: "/results"}
	req.Body = `{"uuid":"0bd74115-2307-458f-8288-b726724045ef"}`
	req.Headers = mock.Values{"Content-Type": {"application/json"}}
	def := &mock.Definition{Request: mock.Request{Path: "/results"}}
	st := Stream{Mock: def, Request: req}

	k := fmt.Sprintf(`file.contents(%s/{{request.body.uuid.regex(\b([0-9a-zA-Z]{4})\b).concat(-x)}}.csv)`, dir)
	result := st.Fill([]string{k})

	if !strings.Contains(result[k][0], "regex chain content") {
		t.Errorf("gjson regex/concat chain was not resolved. Value: %s", result[k])
	}
}

func TestStreamHTTPContentsUnchanged(t *testing.T) {
	// http.contents URLs are kept literal even with request context, to avoid an SSRF
	// surface: a {{request.*}} in the URL must NOT be interpolated.
	req := &mock.Request{Path: "/results/acme"}
	def := &mock.Definition{Request: mock.Request{Path: "/results/:tenant"}}
	st := Stream{Mock: def, Request: req}

	got := st.getInputParam("http.contents(http://example.com/{{request.path.tenant}}/data)", false)
	want := "http://example.com/{{request.path.tenant}}/data"
	if got != want {
		t.Errorf("http.contents URL should stay literal. Got: %q, want: %q", got, want)
	}
}

func TestStreamLiteralPathBackwardCompatible(t *testing.T) {
	// With no request context, holders must stay literal (no panic, no substitution).
	st := Stream{}

	k := "file.contents(/data/{{request.path.x}})"
	result := st.Fill([]string{k})

	if !strings.Contains(result[k][0], "{{request.path.x}}") {
		t.Errorf("Literal path was not preserved. Value: %s", result[k])
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
