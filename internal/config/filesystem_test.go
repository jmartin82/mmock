package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/jmartin82/mmock/v3/internal/config/parser"
	"github.com/jmartin82/mmock/v3/pkg/mock"
)

var MockContent = "{\"URI\":\"test\",\"description\":\"\",\"request\":{\"host\":\"\",\"method\":\"\",\"path\":\"\",\"queryStringParameters\":null,\"headers\":null,\"cookies\":null,\"body\":\"\"},\"response\":{\"statusCode\":0,\"headers\":null,\"cookies\":null,\"body\":\"\"},\"control\":{\"priority\":0,\"delay\":0,\"crazy\":false,\"scenario\":{\"name\":\"\",\"requiredState\":null,\"newState\":\"\"},\"proxyBaseURL\":\"\"}})"

type mockParser struct {
	canParse bool
	readOk   bool
}

func (mr *mockParser) CanParse(filename string) bool {
	return mr.canParse
}

func (mr *mockParser) Parse(content []byte) (mock.Definition, error) {
	if mr.readOk {
		m := mock.Definition{}
		m.URI = "test"
		return m, nil
	}

	return mock.Definition{}, errors.New("error reading")

}

func (mr *mockParser) Marshal(m mock.Definition) ([]byte, error) {
	return []byte("mock"), nil
}

func TestReadMockDefinition(t *testing.T) {

	content := []byte(MockContent)
	dir, err := ioutil.TempDir("", "mmock")
	if err != nil {
		t.Error("Error creating temporary folder")
	}

	tmpfn := filepath.Join(dir, "tmpfile_1")
	if err := os.WriteFile(tmpfn, content, 0666); err != nil {
		t.Error("Error creating temporary file")
	}

	defer os.RemoveAll(dir) // clean up

	mapper := NewFileSystemMapper()
	if _, err := mapper.Read(tmpfn); err == nil {
		t.Error("Expected read error")
	}

	mapper = NewFileSystemMapper()
	mapper.AddParser(&mockParser{canParse: false, readOk: false})
	if _, err := mapper.Read(tmpfn); err == nil {
		t.Error("Expected read error")
	}

	mapper = NewFileSystemMapper()
	mapper.AddParser(&mockParser{canParse: true, readOk: false})
	if _, err := mapper.Read(tmpfn); err == nil {
		t.Error("Expected read error")
	}
	mapper = NewFileSystemMapper()
	mapper.AddParser(&mockParser{canParse: true, readOk: true})
	if _, err := mapper.Read(tmpfn); err != nil {
		t.Errorf("Expected read %v", err)
	}

}

func TestWriteMockDefinition(t *testing.T) {

	mock := mock.Definition{URI: "test"}
	dir, err := ioutil.TempDir("", "mmock")
	if err != nil {
		t.Errorf("Error creating temporary folder")
	}
	defer os.RemoveAll(dir) // clean up

	mapper := NewFileSystemMapper()
	mapper.AddParser(parser.JSONReader{})
	mapper.AddParser(parser.YAMLReader{})

	// Unrecognized extension falls back to JSON.
	tmpfn := filepath.Join(dir, "tmpfile_2")
	if err := mapper.Write(tmpfn, mock); err != nil {
		t.Error("Unexpected error writing the config", err)
	}
	bytes, erf := ioutil.ReadFile(tmpfn)
	if erf != nil || len(bytes) == 0 {
		t.Error("Unexpected error reading the config ", erf)
	}
	if err := json.Unmarshal(bytes, &mock); err != nil {
		t.Error("Expected fallback content to be valid JSON", err)
	}

	// .json extension writes JSON content.
	jsonFn := filepath.Join(dir, "tmpfile_2.json")
	if err := mapper.Write(jsonFn, mock); err != nil {
		t.Error("Unexpected error writing the config", err)
	}
	bytes, erf = ioutil.ReadFile(jsonFn)
	if erf != nil || len(bytes) == 0 {
		t.Error("Unexpected error reading the config ", erf)
	}
	if err := json.Unmarshal(bytes, &mock); err != nil {
		t.Error("Expected .json file to contain valid JSON", err)
	}

	// .yaml extension writes YAML content and must not be overwritten with JSON.
	yamlFn := filepath.Join(dir, "tmpfile_2.yaml")
	if err := mapper.Write(yamlFn, mock); err != nil {
		t.Error("Unexpected error writing the config", err)
	}
	bytes, erf = ioutil.ReadFile(yamlFn)
	if erf != nil || len(bytes) == 0 {
		t.Error("Unexpected error reading the config ", erf)
	}
	if err := yaml.Unmarshal(bytes, &mock); err != nil {
		t.Error("Expected .yaml file to contain valid YAML", err)
	}
	if json.Valid(bytes) {
		t.Error("Expected .yaml file not to be overwritten with JSON content")
	}

}
