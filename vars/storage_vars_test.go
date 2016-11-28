package vars

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/persist"
	"github.com/jmartin82/mmock/vars/fakedata"
)

func getStorageProcessor(persistPath string) VarsProcessor {
	filePersist := persist.NewFilePersister(persistPath)
	persistBag := persist.GetNewPersistEngineBag(filePersist)
	return VarsProcessor{FillerFactory: MockFillerFactory{}, FakeAdapter: fakedata.NewDummyDataFaker("AleixMG"), PersistEngines: persistBag}
}

func TestStorageVars_Sequence_NoParameterPassed(t *testing.T) {
	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)
	processor := getStorageProcessor(persistPath)

	req := definition.Request{}
	res := definition.Response{}
	res.Body = "Sequence: {{storage.Sequence}}"

	mock := definition.Mock{Request: req, Response: res}
	processor.Eval(&req, &mock)

	if mock.Response.Body != "Sequence: {{storage.Sequence}}" {
		t.Error("Replaced tags in body not match", mock.Response.Body)
	}
}

func TestStorageVars_Sequence_SingleParameterPassed(t *testing.T) {
	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)
	processor := getStorageProcessor(persistPath)

	req := definition.Request{}
	res := definition.Response{}

	res.Body = "Sequence: {{storage.Sequence(test)}}"

	mock := definition.Mock{Request: req, Response: res}

	oldValue, _ := processor.PersistEngines.Get("default").GetSequence("test", 0)

	processor.Eval(&req, &mock)

	if mock.Response.Body != "Sequence: "+strconv.Itoa(oldValue) {
		t.Error("Replaced tags in body not match", mock.Response.Body)
	}
}

func TestStorageVars_Sequence_ParameterAndIncreasePassed1(t *testing.T) {
	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)
	processor := getStorageProcessor(persistPath)

	req := definition.Request{}
	res := definition.Response{}

	res.Body = "Sequence: {{storage.Sequence(test, 5)}}"

	mock := definition.Mock{Request: req, Response: res}

	oldValue, _ := processor.PersistEngines.Get("default").GetSequence("test", 0)

	processor.Eval(&req, &mock)

	if mock.Response.Body != "Sequence: "+strconv.Itoa(oldValue+5) {
		t.Error("Replaced tags in body not match", mock.Response.Body)
	}
}

func TestStorageVars_Sequence_ParameterAndIncreasePassed2(t *testing.T) {
	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)
	processor := getStorageProcessor(persistPath)

	req := definition.Request{}
	res := definition.Response{}

	res.Body = "Sequence: {{storage.Sequence('test', 5)}}"

	mock := definition.Mock{Request: req, Response: res}

	oldValue, _ := processor.PersistEngines.Get("default").GetSequence("test", 0)

	processor.Eval(&req, &mock)

	if mock.Response.Body != "Sequence: "+strconv.Itoa(oldValue+5) {
		t.Error("Replaced tags in body not match", mock.Response.Body)
	}
}

func TestStorageVars_GetValue_NoParameterPassed(t *testing.T) {
	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)
	processor := getStorageProcessor(persistPath)

	req := definition.Request{}
	res := definition.Response{}

	res.Body = "GetValue: {{storage.GetValue}}"

	mock := definition.Mock{Request: req, Response: res}
	processor.Eval(&req, &mock)

	if mock.Response.Body != "GetValue: {{storage.GetValue}}" {
		t.Error("Replaced tags in body not match", mock.Response.Body)
	}
}

func TestStorageVars_GetValue_NonExistingKey(t *testing.T) {
	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)
	processor := getStorageProcessor(persistPath)

	req := definition.Request{}
	res := definition.Response{}

	res.Body = "GetValue: {{storage.GetValue(non-existing)}}"

	mock := definition.Mock{Request: req, Response: res}
	processor.Eval(&req, &mock)

	if mock.Response.Body != "GetValue: {{storage.GetValue(non-existing)}}" {
		t.Error("Replaced tags in body not match", mock.Response.Body)
	}
}

func TestStorageVars_GetValue_WithExistingKey(t *testing.T) {
	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)
	processor := getStorageProcessor(persistPath)

	req := definition.Request{}
	res := definition.Response{}

	res.Body = "GetValue: {{storage.GetValue(existing)  }}"

	mock := definition.Mock{Request: req, Response: res}

	processor.PersistEngines.Get("default").SetValue("existing", "test-123")

	processor.Eval(&req, &mock)

	if mock.Response.Body != "GetValue: test-123" {
		t.Error("Replaced tags in body not match", mock.Response.Body)
	}
}

func TestStorageVars_SetValue_NoParametersPassed(t *testing.T) {
	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)
	processor := getStorageProcessor(persistPath)

	req := definition.Request{}
	res := definition.Response{}

	res.Body = "SetValue: {{storage.SetValue}}"

	mock := definition.Mock{Request: req, Response: res}
	processor.Eval(&req, &mock)

	if mock.Response.Body != "SetValue: {{storage.SetValue}}" {
		t.Error("Replaced tags in body not match", mock.Response.Body)
	}
}

func TestStorageVars_SetValue(t *testing.T) {
	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)
	processor := getStorageProcessor(persistPath)

	req := definition.Request{}
	res := definition.Response{}

	res.Body = "SetValue: {{ storage.SetValue(key, test-123) }}"

	mock := definition.Mock{Request: req, Response: res}
	processor.Eval(&req, &mock)

	if mock.Response.Body != "SetValue: test-123" {
		t.Error("Replaced tags in body not match", mock.Response.Body)
	}

	storedValue, err := processor.PersistEngines.Get("default").GetValue("key")
	if err != nil {
		t.Error("Error getting value", err.Error())
	}
	if storedValue != "test-123" {
		t.Error("Stored value not expected", storedValue)
	}
}

func TestStorageVars_SetValue_FromRequestUrl(t *testing.T) {
	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)
	processor := getStorageProcessor(persistPath)

	req := definition.Request{}
	res := definition.Response{}

	req.Path = "/users/123"
	res.Body = "Request+SetValue: {{storage.SetValue(key, {{request.url./users/(?P<value>\\d+)}})}}"

	mock := definition.Mock{Request: req, Response: res}
	processor.Eval(&req, &mock)

	if mock.Response.Body != "Request+SetValue: 123" {
		t.Error("Replaced tags in body not match", mock.Response.Body)
	}
}

func TestStorageVars_SetValue_FromStorage(t *testing.T) {
	persistPath, _ := filepath.Abs("./test_persist")
	defer os.RemoveAll(persistPath)

	os.RemoveAll(persistPath)
	processor := getStorageProcessor(persistPath)

	req := definition.Request{}
	res := definition.Response{}

	res.Body = "Request+SetValue: {{storage.SetValue(key, {{storage.Sequence(tests, 1)}})}}"

	mock := definition.Mock{Request: req, Response: res}
	processor.Eval(&req, &mock)

	if mock.Response.Body != "Request+SetValue: 1" {
		t.Error("Replaced tags in body not match", mock.Response.Body)
	}
}
