package persist

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/parse"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func folderHasFiles(path string) (bool, error) {
	folderExists, _ := exists(path)
	if !folderExists {
		return false, nil
	}
	files, err := ioutil.ReadDir(path)
	return len(files) > 0, err
}

func TestFilePersister_Persist_NoPersistName(t *testing.T) {
	req := definition.Request{}
	res := definition.Response{}
	per := definition.Persist{}

	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)

	persister := NewFilePersister(persistPath, parser)
	persister.Persist(&per, &req, &res)

	hasFiles, _ := folderHasFiles(persistPath)

	if hasFiles {
		t.Error("No file should be created")
	}
}

func TestFilePersister_Persist_FileNotUnderPersistPath(t *testing.T) {
	req := definition.Request{}
	res := definition.Response{}
	per := definition.Persist{}

	res.Body = "BodyToSave"
	per.Name = "../../testing.json"

	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)

	persister := NewFilePersister(persistPath, parser)
	persister.Persist(&per, &req, &res)

	if !strings.HasPrefix(res.Body, "File path not under the persist path.") {
		t.Error("We should end up with an error as the path to the file is not under persist path", res.Body)
	} else if res.StatusCode != 500 {
		t.Error("Status code should be 500", res.StatusCode)
	}
}

func TestFilePersister_Persist_WithBodyToSave(t *testing.T) {
	req := definition.Request{}
	res := definition.Response{}
	per := definition.Persist{}

	res.Body = "BodyToSave"
	per.Name = "testing.json"

	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)

	persister := NewFilePersister(persistPath, parser)
	persister.Persist(&per, &req, &res)

	hasFiles, _ := folderHasFiles(persistPath)

	if !hasFiles {
		t.Error("One file should be created")
	} else {

		filePath := path.Join(persistPath, per.Name)
		fileExists, _ := exists(filePath)

		if !fileExists {
			t.Error("File should exist", filePath)
		} else {

			fileContent, err := ioutil.ReadFile(filePath)

			if err != nil {
				t.Error("The file should be readable", filePath, err)
			} else {

				if string(fileContent) != res.Body {
					t.Error("File content should match result body", string(fileContent), res.Body)
				}
			}
		}
	}
}

func TestFilePersister_LoadBody(t *testing.T) {
	req := definition.Request{}
	res := definition.Response{}

	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)

	res.Persisted = definition.Persisted{Name: "testing_load.json"}
	persister := NewFilePersister(persistPath, parser)

	filePath := path.Join(persistPath, res.Persisted.Name)

	fileContent := "Body to expext"

	err := ioutil.WriteFile(filePath, []byte(fileContent), 0644)
	if err != nil {
		t.Error("File should be written", err)
	} else {
		persister.LoadBody(&req, &res)

		if res.Body != fileContent {
			t.Error("Result body and file content should be the same", res.Body, fileContent)
		}
	}
}

func TestFilePersister_LoadBody_WithAppend(t *testing.T) {
	req := definition.Request{}
	res := definition.Response{}

	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)

	res.Persisted = definition.Persisted{Name: "testing_load.json"}
	res.Persisted.BodyAppend = "Text to append"
	persister := NewFilePersister(persistPath, parser)

	filePath := path.Join(persistPath, res.Persisted.Name)

	fileContent := "Body to expext"

	err := ioutil.WriteFile(filePath, []byte(fileContent), 0644)
	if err != nil {
		t.Error("File should be written", err)
	} else {
		persister.LoadBody(&req, &res)

		if res.Body != fileContent+res.Persisted.BodyAppend {
			t.Error("Result body and file content plus bodyAppend should be the same", res.Body, fileContent, res.Persisted.BodyAppend)
		}
	}
}

func TestFilePersister_LoadBody_FileNotUnderPersistPath(t *testing.T) {
	req := definition.Request{}
	res := definition.Response{}

	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)

	res.Persisted = definition.Persisted{Name: "../../testing_load.json"}
	persister := NewFilePersister(persistPath, parser)

	persister.LoadBody(&req, &res)

	if !strings.HasPrefix(res.Body, "File path not under the persist path.") {
		t.Error("We should end up with an error as the path to the file is not under persist path", res.Body)
	} else if res.StatusCode != 500 {
		t.Error("Status code should be 500", res.StatusCode)
	}
}

func TestFilePersister_LoadBody_NotFound(t *testing.T) {
	req := definition.Request{}
	res := definition.Response{}

	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)

	res.Persisted = definition.Persisted{Name: "testing_load.json"}
	persister := NewFilePersister(persistPath, parser)

	persister.LoadBody(&req, &res)

	if res.Body != "Not Found" {
		t.Error("Result body should be \"Not Found \"", res.Body)
	} else if res.StatusCode != 404 {
		t.Error("Status code should be 404", res.StatusCode)
	}
}

func TestFilePersister_LoadBody_NotFound_CustomTextAndCode(t *testing.T) {
	req := definition.Request{}
	res := definition.Response{}

	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)

	res.Persisted = definition.Persisted{Name: "testing_load.json"}
	persister := NewFilePersister(persistPath, parser)

	res.Persisted.NotFound.StatusCode = 403
	res.Persisted.NotFound.Body = "Really not found"
	res.Persisted.NotFound.BodyAppend = "Appended text"

	persister.LoadBody(&req, &res)

	if res.Body != res.Persisted.NotFound.Body+res.Persisted.NotFound.BodyAppend {
		t.Error("Result body should equal notFound.Body + notFound.BodyAppend", res.Body, res.Persisted.NotFound.Body, res.Persisted.NotFound.BodyAppend)
	} else if res.StatusCode != res.Persisted.NotFound.StatusCode {
		t.Error("Status code should be equal to notFound.StatusCode", res.StatusCode, res.Persisted.NotFound.StatusCode)
	}
}

func TestNewFilePersister(t *testing.T) {
	res := definition.Response{}

	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)

	NewFilePersister(persistPath, parser)

	folderExists, _ := exists(persistPath)

	if !folderExists {
		t.Error("Folder should be created if not existing", res.Body)
	}
}
