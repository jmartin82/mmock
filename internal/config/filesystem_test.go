package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/jmartin82/mmock/pkg/mock"
)

var MockContent = "{\"URI\":\"test\",\"description\":\"\",\"request\":{\"host\":\"\",\"method\":\"\",\"path\":\"\",\"queryStringParameters\":null,\"headers\":null,\"cookies\":null,\"body\":\"\"},\"response\":{\"statusCode\":0,\"headers\":null,\"cookies\":null,\"body\":\"\"},\"control\":{\"priority\":0,\"delay\":0,\"crazy\":false,\"scenario\":{\"name\":\"\",\"requiredState\":null,\"newState\":\"\"},\"proxyBaseURL\":\"\"}})"

type mockParser struct {
	canParse bool
	readOk   bool
}

func (mr *mockParser) CanParse(filename string) bool {
	return mr.canParse
}

func (mr *mockParser) Parse(content []byte) ([]mock.Definition, error) {
	if mr.readOk {
		m := mock.Definition{}
		m.URI = "test"
		return []mock.Definition{m}, nil
	}

	return nil, errors.New("error reading")

}

func TestReadMockDefinition(t *testing.T) {

	content := []byte(MockContent)
	dir, err := ioutil.TempDir("", "mmock")
	if err != nil {
		t.Error("Error creating temporary folder")
	}

	tmpfn := filepath.Join(dir, "tmpfile_1")
	if err := ioutil.WriteFile(tmpfn, content, 0666); err != nil {
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
	tmpfn := filepath.Join(dir, "tmpfile_2")
	defer os.RemoveAll(dir) // clean up
	mapper := NewFileSystemMapper()
	err = mapper.Write(tmpfn, mock)
	if err != nil {
		t.Error("Unexpected error writing the config", err)
	}

	bytes, erf := ioutil.ReadFile(tmpfn)
	if erf != nil || len(bytes) == 0 {
		t.Error("Unexpected error reading the config ", erf)
	}

}
