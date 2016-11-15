package persist

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"log"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/parse"
	"gopkg.in/mgo.v2"
)

const mongoTestURL = "mongodb://localhost/mmock_test"

type ConnectionState string

const (
	unknown ConnectionState = "Unknown"
	success ConnectionState = "Success"
	fail    ConnectionState = "Fail"
)

var connectionState = unknown

// using this method to avoid testing if you have no mongo installed on "mongodb://localhost
func hasConnection(t *testing.T) bool {
	switch connectionState {
	case fail:
		return false
	case success:
		return true
	case unknown:
		fallthrough
	default:

		dialInfo, err := mgo.ParseURL(mongoTestURL)
		if err != nil {
			t.Error(err)
			connectionState = fail
			return false
		}

		if dialInfo.Timeout == 0 {
			dialInfo.Timeout = 2 * time.Second
		}

		session, err := mgo.DialWithInfo(dialInfo)
		if err != nil {
			log.Printf("Cannot connect to mongo, make sure you have installed mongo server listening on mongodb://localhost. Error: %s", err.Error())
			connectionState = fail
			return false
		}
		defer session.Close()

		connectionState = success
		return true
	}
}

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
	if !hasConnection(t) {
		return
	}

	req := definition.Request{}
	res := definition.Response{}
	per := definition.Persist{}

	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

	persistPath := mongoTestURL
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
	if !hasConnection(t) {
		return
	}
	req := definition.Request{}
	res := definition.Response{}
	per := definition.Persist{}

	res.Body = "{ \"test\": 1}"
	per.Name = "test/testing-1"

	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

	persistPath := mongoTestURL
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

func TestMongoBodyPersister_LoadBody(t *testing.T) {
	if !hasConnection(t) {
		return
	}
	req := definition.Request{}
	res := definition.Response{}

	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

	persistPath := mongoTestURL
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

	res.Persisted = definition.Persisted{Name: "test/item-1"}

	content := "{\"test\": \"body to expect\"}"
	id := "item-1"

	err = persister.SaveItem("test", id, content)
	if err != nil {
		t.Error("Item should be saved", err)
	} else {
		persister.LoadBody(&req, &res)

		if equal, err := jsonsAreEqual(res.Body, content); !equal || err != nil {
			t.Error("Result body and file content should be the same", res.Body, content)
		}
	}
}

func TestMongoBodyPersister_LoadBody_WithAppend(t *testing.T) {
	if !hasConnection(t) {
		return
	}
	req := definition.Request{}
	res := definition.Response{}

	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

	persistPath := mongoTestURL
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

	res.Persisted = definition.Persisted{Name: "test/item-1", BodyAppend: "{\"append\":1}"}

	content := "{\"test\": \"body to expect\"}"
	id := "item-1"

	err = persister.SaveItem("test", id, content)
	if err != nil {
		t.Error("Item should be saved", err)
	} else {
		persister.LoadBody(&req, &res)

		if equal, err := jsonsAreEqual(res.Body, parse.JoinJSON(content, res.Persisted.BodyAppend)); !equal || err != nil {
			t.Error("Result body and file content plus bodyAppend should be the same", res.Body, content, res.Persisted.BodyAppend)
		}
	}
}

func TestMongoBodyPersister_LoadBody_NotFound(t *testing.T) {
	if !hasConnection(t) {
		return
	}
	req := definition.Request{}
	res := definition.Response{}

	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

	persistPath := mongoTestURL
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

	res.Persisted = definition.Persisted{Name: "test/item-1"}

	persister.LoadBody(&req, &res)

	if res.Body != "Not Found" {
		t.Error("Result body should be \"Not Found \"", res.Body)
	} else if res.StatusCode != 404 {
		t.Error("Status code should be 404", res.StatusCode)
	}
}

func TestMongoBodyPersister_LoadBody_NotFound_CustomTextAndCode(t *testing.T) {
	if !hasConnection(t) {
		return
	}
	req := definition.Request{}
	res := definition.Response{}

	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

	persistPath := mongoTestURL
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

	res.Persisted = definition.Persisted{Name: "test/item-1"}

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
