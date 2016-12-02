package vars

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/persist"
	"github.com/jmartin82/mmock/vars/fakedata"
)

func getFileProcessor(persistPath string) VarsProcessor {
	filePersist := persist.NewFilePersister(persistPath)
	persistBag := persist.GetNewPersistEngineBag(filePersist)
	return VarsProcessor{FillerFactory: MockFillerFactory{}, FakeAdapter: fakedata.NewDummyDataFaker("AleixMG"), PersistEngines: persistBag}
}

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

	actions := make(map[string]string)
	actions["write"] = "{{ request.body }}"
	per := definition.Persist{Actions: actions}

	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)

	varsProcessor := getFileProcessor(persistPath)
	mock := definition.Mock{Request: req, Response: res, Persist: per}
	varsProcessor.Eval(&req, &mock)

	hasFiles, _ := folderHasFiles(persistPath)

	if hasFiles {
		t.Error("No file should be created")
	}
}

func TestFilePersister_Persist_FileNotUnderPersistPath(t *testing.T) {
	req := definition.Request{}
	res := definition.Response{Body: "{{persist.entity.content}}", StatusCode: 200}

	actions := make(map[string]string)
	actions["write"] = "{{ request.body }}"
	per := definition.Persist{Actions: actions}

	req.Body = "BodyToSave"
	per.Entity = "../../testing.json"

	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)

	varsProcessor := getFileProcessor(persistPath)
	mock := definition.Mock{Request: req, Response: res, Persist: per}
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != "" {
		t.Error("We should end up with an error as the path to the file is not under persist path", mock.Response.Body)
	} else if mock.Response.StatusCode != 404 {
		t.Error("Status code should be 404", mock.Response.StatusCode)
	}
}

func TestFilePersister_Persist_WithBodyToSave(t *testing.T) {
	req := definition.Request{}
	res := definition.Response{Body: "{{persist.entity.content}}", StatusCode: 200}

	actions := make(map[string]string)
	actions["write"] = "{{ request.body }}"
	per := definition.Persist{Actions: actions}

	req.Body = "BodyToSave"
	per.Entity = "testing.json"

	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)

	varsProcessor := getFileProcessor(persistPath)
	mock := definition.Mock{Request: req, Response: res, Persist: per}
	varsProcessor.Eval(&req, &mock)

	hasFiles, _ := folderHasFiles(persistPath)

	if !hasFiles {
		t.Error("One file should be created")
	} else {

		filePath := path.Join(persistPath, mock.Persist.Entity)
		fileExists, _ := exists(filePath)

		if !fileExists {
			t.Error("File should exist", filePath)
		} else {

			fileContent, err := ioutil.ReadFile(filePath)

			if err != nil {
				t.Error("The file should be readable", filePath, err)
			} else {

				stringContent := string(fileContent)

				if stringContent != mock.Request.Body {
					t.Error("File content should match result body", stringContent, mock.Response.Body)
				}
			}
		}
	}
}

func TestFilePersister_LoadBody(t *testing.T) {
	req := definition.Request{}
	res := definition.Response{Body: "{{ persist.entity.content }}"}
	per := definition.Persist{Entity: "testing_load.json"}

	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)

	varsProcessor := getFileProcessor(persistPath)
	mock := definition.Mock{Request: req, Response: res, Persist: per}

	filePath := path.Join(persistPath, per.Entity)

	fileContent := "Body to expect"

	err := ioutil.WriteFile(filePath, []byte(fileContent), 0755)
	if err != nil {
		t.Error("File should be written", err)
	} else {
		varsProcessor.Eval(&req, &mock)

		if mock.Response.Body != fileContent {
			t.Error("Result body and file content should be the same", mock.Response.Body, fileContent)
		}
	}
}

func TestFilePersister_LoadBody_WithAppend(t *testing.T) {
	req := definition.Request{}
	res := definition.Response{Body: "{{persist.entity.content}}"}

	appendText := "Text to append"

	actions := make(map[string]string)
	actions["append"] = appendText

	per := definition.Persist{Entity: "testing_load.json", Actions: actions}

	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)

	varsProcessor := getFileProcessor(persistPath)
	mock := definition.Mock{Request: req, Response: res, Persist: per}

	filePath := path.Join(persistPath, per.Entity)

	fileContent := "Body to expect"

	err := ioutil.WriteFile(filePath, []byte(fileContent), 0755)
	if err != nil {
		t.Error("File should be written", err)
	} else {
		varsProcessor.Eval(&req, &mock)

		if mock.Response.Body != fileContent+appendText {
			t.Error("Result body and file content plus bodyAppend should be the same", mock.Response.Body, fileContent, appendText)
		}
	}
}

func TestFilePersister_LoadBody_FileNotUnderPersistPath(t *testing.T) {
	req := definition.Request{}
	res := definition.Response{Body: "{{persist.entity.content}}"}
	per := definition.Persist{Entity: "../../testing_load.json"}

	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)

	varsProcessor := getFileProcessor(persistPath)
	mock := definition.Mock{Request: req, Response: res, Persist: per}
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != "" {
		t.Error("We should end up with an empty body as the path to the file is not under persist path", mock.Response.Body)
	} else if mock.Response.StatusCode != 404 {
		t.Error("Status code should be 404", mock.Response.StatusCode)
	}
}

func TestFilePersister_LoadBody_NotFound(t *testing.T) {
	req := definition.Request{}
	res := definition.Response{Body: "{{persist.entity.content}}"}
	per := definition.Persist{Entity: "testing_load.json"}

	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)

	varsProcessor := getFileProcessor(persistPath)
	mock := definition.Mock{Request: req, Response: res, Persist: per}
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != "" {
		t.Error("Result body should be empty", mock.Response.Body)
	} else if mock.Response.StatusCode != 404 {
		t.Error("Status code should be 404", mock.Response.StatusCode)
	}
}

func TestNewFilePersister(t *testing.T) {
	res := definition.Response{}

	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)

	persist.NewFilePersister(persistPath)

	folderExists, _ := exists(persistPath)

	if !folderExists {
		t.Error("Folder should be created if not existing", res.Body)
	}
}


