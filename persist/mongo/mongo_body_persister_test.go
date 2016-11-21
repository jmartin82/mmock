package mongo

// import (
// 	"log"
// 	"testing"
// 	"time"

// 	"github.com/jmartin82/mmock/definition"
// 	"github.com/jmartin82/mmock/parse"
// 	"gopkg.in/mgo.v2"
// )

// const mongoTestURL = "mongodb://localhost/mmock_test"

// type ConnectionState string

// const (
// 	unknown ConnectionState = "Unknown"
// 	success ConnectionState = "Success"
// 	fail    ConnectionState = "Fail"
// )

// var connectionState = unknown

// // using this method to avoid testing if you have no mongo installed on "mongodb://localhost
// func hasConnection(t *testing.T) bool {
// 	switch connectionState {
// 	case fail:
// 		return false
// 	case success:
// 		return true
// 	case unknown:
// 		fallthrough
// 	default:

// 		dialInfo, err := mgo.ParseURL(mongoTestURL)
// 		if err != nil {
// 			t.Error(err)
// 			connectionState = fail
// 			return false
// 		}

// 		if dialInfo.Timeout == 0 {
// 			dialInfo.Timeout = 2 * time.Second
// 		}

// 		session, err := mgo.DialWithInfo(dialInfo)
// 		if err != nil {
// 			log.Printf("Cannot connect to mongo, make sure you have installed mongo server listening on mongodb://localhost. Error: %s", err.Error())
// 			connectionState = fail
// 			return false
// 		}
// 		defer session.Close()

// 		connectionState = success
// 		return true
// 	}
// }

// func mongoCollectionHasItems(collection *mgo.Collection) (bool, error) {
// 	count, err := collection.Count()
// 	return count > 0, err
// }

// func TestMongoBodyPersister_Persist_NoPersistName(t *testing.T) {
// 	if !hasConnection(t) {
// 		return
// 	}

// 	req := definition.Request{}
// 	res := definition.Response{}
// 	per := definition.Persist{}

// 	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

// 	persister := NewMongoBodyPersister(mongoTestURL, parser)
// 	defer func() {
// 		persister.Repository.dropDatabase() // cleanup database
// 	}()
// 	persister.Repository.dropDatabase() // make sure we are working on a clean database

// 	persister.Persist(&per, &req, &res)

// 	hasCollections, err := persister.Repository.hasCollections()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if hasCollections {
// 		t.Error("No collections should be present")
// 	}
// }

// func TestMongoBodyPersister_Persist_WithBodyToSave(t *testing.T) {
// 	if !hasConnection(t) {
// 		return
// 	}
// 	req := definition.Request{}
// 	res := definition.Response{}
// 	per := definition.Persist{}

// 	res.Body = "{ \"test\": 1}"
// 	per.Name = "test/testing-1"

// 	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

// 	persister := NewMongoBodyPersister(mongoTestURL, parser)
// 	defer func() {
// 		persister.Repository.dropDatabase() // cleanup database
// 	}()
// 	persister.Repository.dropDatabase() // make sure we are working on a clean database

// 	persister.Persist(&per, &req, &res)

// 	hasItems, _ := persister.Repository.hasCollectionsItems("test")

// 	if !hasItems {
// 		t.Error("One item should be created")
// 	} else {

// 		bodyContent, err := persister.Repository.GetItem("test", "testing-1")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		if equal, err := parse.JSONSStringsAreEqual(bodyContent, res.Body); !equal || err != nil {
// 			t.Error("File content should match result body", bodyContent, res.Body)
// 		}
// 	}
// }

// func TestMongoBodyPersister_Persist_WithBodyToSave_WrongName(t *testing.T) {
// 	if !hasConnection(t) {
// 		return
// 	}
// 	req := definition.Request{}
// 	res := definition.Response{}
// 	per := definition.Persist{}

// 	res.Body = "{ \"test\": 1}"
// 	per.Name = "test-testing-1"

// 	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

// 	persister := NewMongoBodyPersister(mongoTestURL, parser)
// 	defer func() {
// 		persister.Repository.dropDatabase() // cleanup database
// 	}()
// 	persister.Repository.dropDatabase() // make sure we are working on a clean database

// 	persister.Persist(&per, &req, &res)

// 	if res.StatusCode != 500 {
// 		t.Error("The status code should be 500 as the persist name is in wrong format. Current status code is: ", res.StatusCode)
// 	}
// 	if res.Body != errWrongNameFormat.Error() {
// 		t.Error("The body should contain the error message for wrong name format. Current body content is: ", res.Body)
// 	}
// }

// func TestMongoBodyPersister_Persist_WithNonJSONBodyToSave(t *testing.T) {
// 	if !hasConnection(t) {
// 		return
// 	}
// 	req := definition.Request{}
// 	res := definition.Response{}
// 	per := definition.Persist{}

// 	res.Body = "Body to save"
// 	per.Name = "test/testing-1"

// 	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

// 	persister := NewMongoBodyPersister(mongoTestURL, parser)
// 	defer func() {
// 		persister.Repository.dropDatabase() // cleanup database
// 	}()
// 	persister.Repository.dropDatabase() // make sure we are working on a clean database

// 	persister.Persist(&per, &req, &res)

// 	hasItems, _ := persister.Repository.hasCollectionsItems("test")

// 	if !hasItems {
// 		t.Error("One item should be created")
// 	} else {

// 		bodyContent, err := persister.Repository.GetItem("test", "testing-1")
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		if equal, err := parse.JSONSStringsAreEqual(bodyContent, res.Body); !equal || err != nil {
// 			t.Error("File content should match result body", bodyContent, res.Body)
// 		}
// 	}
// }

// func TestMongoBodyPersister_LoadBody(t *testing.T) {
// 	if !hasConnection(t) {
// 		return
// 	}
// 	req := definition.Request{}
// 	res := definition.Response{}

// 	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

// 	persister := NewMongoBodyPersister(mongoTestURL, parser)
// 	defer func() {
// 		persister.Repository.dropDatabase() // cleanup database
// 	}()
// 	persister.Repository.dropDatabase() // make sure we are working on a clean database

// 	res.Persisted = definition.Persisted{Name: "test/item-1"}

// 	content := "{\"test\": \"body to expect\"}"
// 	id := "item-1"

// 	err := persister.Repository.UpsertItem("test", id, content)
// 	if err != nil {
// 		t.Error("Item should be saved", err)
// 	} else {
// 		persister.LoadBody(&req, &res)

// 		if equal, err := parse.JSONSStringsAreEqual(res.Body, content); !equal || err != nil {
// 			t.Error("Result body and file content should be the same", res.Body, content)
// 		}
// 	}
// }

// func TestMongoBodyPersister_LoadBody_WrongName(t *testing.T) {
// 	if !hasConnection(t) {
// 		return
// 	}
// 	req := definition.Request{}
// 	res := definition.Response{}

// 	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

// 	persister := NewMongoBodyPersister(mongoTestURL, parser)
// 	defer func() {
// 		persister.Repository.dropDatabase() // cleanup database
// 	}()
// 	persister.Repository.dropDatabase() // make sure we are working on a clean database

// 	res.Persisted = definition.Persisted{Name: "test-item-1"}

// 	persister.LoadBody(&req, &res)

// 	if res.StatusCode != 500 {
// 		t.Error("The status code should be 500 as the persist name is in wrong format. Current status code is: ", res.StatusCode)
// 	}
// 	if res.Body != errWrongNameFormat.Error() {
// 		t.Error("The body should contain the error message for wrong name format. Current body content is: ", res.Body)
// 	}
// }

// func TestMongoBodyPersister_LoadBodyNonJSON(t *testing.T) {
// 	if !hasConnection(t) {
// 		return
// 	}
// 	req := definition.Request{}
// 	res := definition.Response{}

// 	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

// 	persister := NewMongoBodyPersister(mongoTestURL, parser)
// 	defer func() {
// 		persister.Repository.dropDatabase() // cleanup database
// 	}()
// 	persister.Repository.dropDatabase() // make sure we are working on a clean database

// 	res.Persisted = definition.Persisted{Name: "test/item-1"}

// 	content := "Non JSON body to expect"
// 	id := "item-1"

// 	err := persister.Repository.UpsertItem("test", id, content)
// 	if err != nil {
// 		t.Error("Item should be saved", err)
// 	} else {
// 		persister.LoadBody(&req, &res)

// 		if equal, err := parse.JSONSStringsAreEqual(res.Body, content); !equal || err != nil {
// 			t.Error("Result body and file content should be the same", res.Body, content)
// 		}
// 	}
// }

// func TestMongoBodyPersister_LoadBody_WithAppend(t *testing.T) {
// 	if !hasConnection(t) {
// 		return
// 	}
// 	req := definition.Request{}
// 	res := definition.Response{}

// 	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

// 	persister := NewMongoBodyPersister(mongoTestURL, parser)
// 	defer func() {
// 		persister.Repository.dropDatabase() // cleanup database
// 	}()
// 	persister.Repository.dropDatabase() // make sure we are working on a clean database

// 	res.Persisted = definition.Persisted{Name: "test/item-1", BodyAppend: "{\"append\":1}"}

// 	content := "{\"test\": \"body to expect\"}"
// 	id := "item-1"

// 	err := persister.Repository.UpsertItem("test", id, content)
// 	if err != nil {
// 		t.Error("Item should be saved", err)
// 	} else {
// 		persister.LoadBody(&req, &res)

// 		if equal, err := parse.JSONSStringsAreEqual(res.Body, parse.JoinJSON(content, res.Persisted.BodyAppend)); !equal || err != nil {
// 			t.Error("Result body and file content plus bodyAppend should be the same", res.Body, content, res.Persisted.BodyAppend)
// 		}
// 	}
// }

// func TestMongoBodyPersister_LoadBody_NotFound(t *testing.T) {
// 	if !hasConnection(t) {
// 		return
// 	}
// 	req := definition.Request{}
// 	res := definition.Response{}

// 	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

// 	persister := NewMongoBodyPersister(mongoTestURL, parser)
// 	defer func() {
// 		persister.Repository.dropDatabase() // cleanup database
// 	}()
// 	persister.Repository.dropDatabase() // make sure we are working on a clean database

// 	res.Persisted = definition.Persisted{Name: "test/item-1"}

// 	persister.LoadBody(&req, &res)

// 	if res.Body != "Not Found" {
// 		t.Error("Result body should be \"Not Found \"", res.Body)
// 	} else if res.StatusCode != 404 {
// 		t.Error("Status code should be 404", res.StatusCode)
// 	}
// }

// func TestMongoBodyPersister_LoadBody_NotFound_CustomTextAndCode(t *testing.T) {
// 	if !hasConnection(t) {
// 		return
// 	}
// 	req := definition.Request{}
// 	res := definition.Response{}

// 	parser := parse.FakeDataParse{Fake: parse.DummyDataFaker{Dummy: "AleixMG"}}

// 	persister := NewMongoBodyPersister(mongoTestURL, parser)
// 	defer func() {
// 		persister.Repository.dropDatabase() // cleanup database
// 	}()
// 	persister.Repository.dropDatabase() // make sure we are working on a clean database

// 	res.Persisted = definition.Persisted{Name: "test/item-1"}

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
