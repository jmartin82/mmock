package config

import (
	"github.com/jmartin82/mmock/v3/pkg/mock"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type DummyReader struct {
}

//CanParse return true if is a json file
func (jp DummyReader) CanParse(filename string) bool {

	return strings.Contains(filename, "tmpfile_populate")
}

//Read Unmarshal a json file to Mock struct
func (jp DummyReader) Parse(buf []byte) (mock.Definition, error) {
	return mock.Definition{URI: "test_mapping"}, nil
}

var fsUpdate = make(chan struct{})

func TestMappingCrud(t *testing.T) {
	dir, err := ioutil.TempDir("", "mmock_mapping")
	if err != nil {
		t.Errorf("Error creating temporary folder")
	}
	defer os.RemoveAll(dir)
	cm := NewFileSystemMapper()
	n := NewConfigMapping(dir, cm, fsUpdate)

	if _, f := n.Get("mock_definition_example"); f {
		t.Errorf("Error mapping shouldn't exist")
	}

	if err := n.Set("mock_definition_example", mock.Definition{URI: "test"}); err != nil {
		t.Errorf("Error persisting mapping")
	}

	if m, f := n.Get("mock_definition_example"); !f || m.URI != "test" {
		t.Errorf("Error mapping should exist")
	}

	if err := n.Set("mock_definition_example2", mock.Definition{URI: "test2"}); err != nil {
		t.Errorf("Error persisting mapping")
	}

	mocks := n.List()
	if len(mocks) != 2 {
		t.Errorf("Error listing mapping")
	}

	if err := n.Delete("mock_definition_example"); err != nil {
		t.Errorf("Error deletegin mapping")
	}

	if _, f := n.Get("mock_definition_example"); f {
		t.Errorf("Error mapping shouldn't exist")
	}

	mocks = n.List()
	if len(mocks) != 1 {
		t.Errorf("Error listing mapping")
	}

}

func TestPopulate(t *testing.T) {
	dir, err := ioutil.TempDir("", "mmock_mapping")
	if err != nil {
		t.Errorf("Error creating temporary folder")
	}

	defer os.RemoveAll(dir) // clean up

	cm := NewFileSystemMapper()
	cm.AddParser(DummyReader{})
	n := NewConfigMapping(dir, cm, fsUpdate)
	mocks := n.List()
	if len(mocks) != 0 {
		t.Errorf("Error listing mapping, should be 0")
	}

	tmpfn := filepath.Join(dir, "tmpfile_populate")
	if err := ioutil.WriteFile(tmpfn, []byte("some content"), 0666); err != nil {
		t.Errorf("Error creating temporary file")
	}

	n.populate()

	mocks = n.List()
	if len(mocks) != 1 {
		t.Errorf("Error listing mapping should be 1")
	}

	if m, f := n.Get("tmpfile_populate"); !f || m.URI != "tmpfile_populate" {
		t.Errorf("Error mapping should exist")
	}

}

func TestGetSortedMappingList(t *testing.T) {

	dir, err := ioutil.TempDir("", "mmock_mapping")
	if err != nil {
		t.Errorf("Error creating temporary folder")
	}
	defer os.RemoveAll(dir)
	cm := NewFileSystemMapper()
	n := NewConfigMapping(dir, cm, fsUpdate)

	n.Set("mock_definition_example", mock.Definition{URI: "test"})
	n.Set("mock_definition_example2", mock.Definition{URI: "test2", Control: mock.Control{Priority: 1}})
	n.Set("mock_definition_example3", mock.Definition{URI: "test3", Control: mock.Control{Priority: 100}})

	mocks := n.List()

	if mocks[0].URI != "test3" {
		t.Errorf("Invalid priority sort")
	}

	if mocks[1].URI != "test2" {
		t.Errorf("Invalid priority sort")
	}

	if mocks[2].URI != "test" {
		t.Errorf("Invalid priority sort")
	}

}
