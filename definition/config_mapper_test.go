package definition

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var MockContent = "{\"URI\":\"test\",\"description\":\"\",\"request\":{\"host\":\"\",\"method\":\"\",\"path\":\"\",\"queryStringParameters\":null,\"headers\":null,\"cookies\":null,\"body\":\"\"},\"response\":{\"statusCode\":0,\"headers\":null,\"cookies\":null,\"body\":\"\"},\"control\":{\"priority\":0,\"delay\":0,\"crazy\":false,\"scenario\":{\"name\":\"\",\"requiredState\":null,\"newState\":\"\"},\"proxyBaseURL\":\"\"}})"

type mockParser struct {
	canParse bool
	readOk   bool
}

func (mr *mockParser) CanParse(filename string) bool {
	return mr.canParse
}

func (mr *mockParser) Parse(content []byte) (Mock, error) {
	if mr.readOk {
		m := Mock{}
		m.URI = "test"
		return m, nil
	}

	return Mock{}, errors.New("error reading")

}

func TestReadMockDefinition(t *testing.T) {

	content := []byte(MockContent)
	dir, err := ioutil.TempDir("", "mmock")
	if err != nil {
		t.Errorf("Error creating temporary folder")
	}

	tmpfn := filepath.Join(dir, "tmpfile_1")
	if err := ioutil.WriteFile(tmpfn, content, 0666); err != nil {
		t.Errorf("Error creating temporary file")
	}

	defer os.RemoveAll(dir) // clean up

	mapper := NewConfigMapper()
	if _, err := mapper.Read(tmpfn); err == nil {
		t.Errorf("Expected read error")
	}

	mapper = NewConfigMapper()
	mapper.AddConfigParser(&mockParser{canParse: false, readOk: false})
	if _, err := mapper.Read(tmpfn); err == nil {
		t.Errorf("Expected read error")
	}

	mapper = NewConfigMapper()
	mapper.AddConfigParser(&mockParser{canParse: true, readOk: false})
	if _, err := mapper.Read(tmpfn); err == nil {
		t.Errorf("Expected read error")
	}
	mapper = NewConfigMapper()
	mapper.AddConfigParser(&mockParser{canParse: true, readOk: true})
	if _, err := mapper.Read(tmpfn); err != nil {
		t.Errorf("Expected read %v", err)
	}

}

func TestWriteMockDefinition(t *testing.T) {

	mock := Mock{URI: "test"}
	dir, err := ioutil.TempDir("", "mmock")
	if err != nil {
		t.Errorf("Error creating temporary folder")
	}
	tmpfn := filepath.Join(dir, "tmpfile_2")
	defer os.RemoveAll(dir) // clean up
	mapper := NewConfigMapper()
	err = mapper.Write(tmpfn, mock)
	if err != nil {
		t.Errorf("Unexpected error writing the config", err)
	}

	bytes, erf := ioutil.ReadFile(tmpfn)
	if erf != nil || len(bytes) == 0 {
		t.Errorf("Unexpected error reading the config ", erf)
	}

}
