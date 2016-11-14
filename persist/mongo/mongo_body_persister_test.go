package persist

import (
	"encoding/json"
	"os"
	"testing"

	mgo "gopkg.in/mgo.v2"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/parse"
)

func formatJSON(input string) (result string, err error) {
	var jsonParsed interface{}
	json.Unmarshal([]byte(input), &jsonParsed)
	if err != nil {
		return "", err
	}

	byteString, err := json.Marshal(jsonParsed)
	if err != nil {
		return "", err
	}

	return string(byteString), nil
}

func jsonsAreEqual(input1 string, input2 string) (result bool, err error) {
	formatedInput1, err := formatJSON(input1)
	if err != nil {
		return false, err
	}
	formatedInput2, err := formatJSON(input2)
	if err != nil {
		return false, err
	}
	return formatedInput1 == formatedInput2, nil
}

func mongoRecordExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func mongoCollectionHasItems(collection *mgo.Collection) (bool, error) {
	count, err := collection.Count()
	return count > 0, err
}

func TestMongoBodyPersister_Persist_NoPersistName(t *testing.T) {
	req := definition.Request{}
	res := definition.Response{}
	per := definition.Persist{}

	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

	persistPath := "mongodb://localhost/mmock_test"
	persister := NewMongoBodyPersister(persistPath, parser)

	session, err := persister.ConnectMongo()
	if err != nil {
		t.Error(err.Error())
		return
	}

	defer func() {
		session.DB(persister.ConnectionInfo.Database).DropDatabase()
		session.Close()
	}()

	persister.Persist(&per, &req, &res)

	collectionNames, _ := session.DB(persister.ConnectionInfo.Database).CollectionNames()

	if len(collectionNames) > 0 {
		t.Error("No collections should be present")
	}
}

func TestMongoBodyPersister_Persist_WithBodyToSave(t *testing.T) {
	req := definition.Request{}
	res := definition.Response{}
	per := definition.Persist{}

	res.Body = "{ \"test\": 1}"
	per.Name = "test/testing-1"

	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

	persistPath := "mongodb://localhost/mmock_test"
	persister := NewMongoBodyPersister(persistPath, parser)

	session, err := persister.ConnectMongo()
	if err != nil {
		t.Error(err.Error())
		return
	}

	defer func() {
		session.DB(persister.ConnectionInfo.Database).DropDatabase()
		session.Close()
	}()

	persister.Persist(&per, &req, &res)

	collection := session.DB(persister.ConnectionInfo.Database).C("test")

	hasItems, _ := mongoCollectionHasItems(collection)

	if !hasItems {
		t.Error("One item should be created")
	} else {

		var result interface{}
		err = collection.FindId("testing-1").One(&result)
		if err != nil {
			t.Error(err.Error())
		} else {

			bodyContent, err := persister.GetResultString(result)
			if err != nil {
				t.Error(err.Error())
			} else {

				if equal, err := jsonsAreEqual(string(bodyContent), res.Body); !equal || err != nil {
					t.Error("File content should match result body", string(bodyContent), res.Body)
				}
			}
		}
	}
}

// func TestMongoBodyPersister_LoadBody(t *testing.T) {
// 	req := definition.Request{}
// 	res := definition.Response{}

// 	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

// 	persistPath := "mongodb://localhost/mmock_test"
// 	defer os.RemoveAll(persistPath)

// 	os.RemoveAll(persistPath)

// 	res.Persisted = definition.Persisted{Name: "testing_load.json"}
// 	persister := NewMongoBodyPersister(persistPath, parser)

// 	filePath := path.Join(persistPath, res.Persisted.Name)

// 	fileContent := "Body to expext"

// 	err := ioutil.WriteFile(filePath, []byte(fileContent), 0644)
// 	if err != nil {
// 		t.Error("File should be written", err)
// 	} else {
// 		persister.LoadBody(&req, &res)

// 		if res.Body != fileContent {
// 			t.Error("Result body and file content should be the same", res.Body, fileContent)
// 		}
// 	}
// }

// func TestMongoBodyPersister_LoadBody_WithAppend(t *testing.T) {
// 	req := definition.Request{}
// 	res := definition.Response{}

// 	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

// 	persistPath := "mongodb://localhost/mmock_test"
// 	defer os.RemoveAll(persistPath)

// 	os.RemoveAll(persistPath)

// 	res.Persisted = definition.Persisted{Name: "testing_load.json"}
// 	res.Persisted.BodyAppend = "Text to append"
// 	persister := NewMongoBodyPersister(persistPath, parser)

// 	filePath := path.Join(persistPath, res.Persisted.Name)

// 	fileContent := "Body to expext"

// 	err := ioutil.WriteFile(filePath, []byte(fileContent), 0644)
// 	if err != nil {
// 		t.Error("File should be written", err)
// 	} else {
// 		persister.LoadBody(&req, &res)

// 		if res.Body != fileContent+res.Persisted.BodyAppend {
// 			t.Error("Result body and file content plus bodyAppend should be the same", res.Body, fileContent, res.Persisted.BodyAppend)
// 		}
// 	}
// }

// func TestMongoBodyPersister_LoadBody_NotFound(t *testing.T) {
// 	req := definition.Request{}
// 	res := definition.Response{}

// 	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

// 	persistPath := "mongodb://localhost/mmock_test"
// 	defer os.RemoveAll(persistPath)

// 	os.RemoveAll(persistPath)

// 	res.Persisted = definition.Persisted{Name: "testing_load.json"}
// 	persister := NewMongoBodyPersister(persistPath, parser)

// 	persister.LoadBody(&req, &res)

// 	if res.Body != "Not Found" {
// 		t.Error("Result body should be \"Not Found \"", res.Body)
// 	} else if res.StatusCode != 404 {
// 		t.Error("Status code should be 404", res.StatusCode)
// 	}
// }

// func TestMongoBodyPersister_LoadBody_NotFound_CustomTextAndCode(t *testing.T) {
// 	req := definition.Request{}
// 	res := definition.Response{}

// 	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

// 	persistPath := "mongodb://localhost/mmock_test"
// 	defer os.RemoveAll(persistPath)

// 	os.RemoveAll(persistPath)

// 	res.Persisted = definition.Persisted{Name: "testing_load.json"}
// 	persister := NewMongoBodyPersister(persistPath, parser)

// 	res.Persisted.NotFound.StatusCode = 403
// 	res.Persisted.NotFound.Body = "Really not found"
// 	res.Persisted.NotFound.BodyAppend = "Appended text"

// 	persister.LoadBody(&req, &res)

// 	if res.Body != res.Persisted.NotFound.Body+res.Persisted.NotFound.BodyAppend {
// 		t.Error("Result body should equal notFound.Body + notFound.BodyAppend", res.Body, res.Persisted.NotFound.Body, res.Persisted.NotFound.BodyAppend)
// 	} else if res.StatusCode != res.Persisted.NotFound.StatusCode {
// 		t.Error("Status code should be equal to notFound.StatusCode", res.StatusCode, res.Persisted.NotFound.StatusCode)
// 	}
// }
