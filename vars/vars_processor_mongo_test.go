package vars

import (
	"log"
	"testing"
	"time"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/persist"
	"github.com/jmartin82/mmock/utils"
	"github.com/jmartin82/mmock/vars/fakedata"
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

func getMongoProcessor() (VarsProcessor, *persist.MongoPersister) {
	mongoPersist := persist.NewMongoPersister(mongoTestURL)
	persistBag := persist.GetNewPersistEngineBag(mongoPersist)
	return VarsProcessor{FillerFactory: MockFillerFactory{}, FakeAdapter: fakedata.NewDummyDataFaker("AleixMG"), PersistEngines: persistBag}, mongoPersist
}

//dropDatabase drops the connected database - TO BE USED IN TESTS ONLY!!!
func dropDatabase(mr persist.MongoRepository) error {

	session, err := mr.ConnectMongo()
	if err != nil {
		return err
	}

	defer session.Close()

	session.DB(mr.ConnectionInfo.Database).DropDatabase()
	return nil
}

//hasCollections returns whether the database has any collections inside or not
func hasCollections(mr persist.MongoRepository) (bool, error) {

	session, err := mr.ConnectMongo()
	if err != nil {
		return false, err
	}

	defer session.Close()

	collectionNames, err := session.DB(mr.ConnectionInfo.Database).CollectionNames()
	return len(collectionNames) > 0, err
}

//hasCollectionsItems returns whether the a given collection has items
func hasCollectionsItems(mr persist.MongoRepository, collectionName string) (bool, error) {

	session, err := mr.ConnectMongo()
	if err != nil {
		return false, err
	}

	defer session.Close()

	collection := session.DB(mr.ConnectionInfo.Database).C(collectionName)
	count, err := collection.Count()
	return count > 0, err
}

func mongoCollectionHasItems(collection *mgo.Collection) (bool, error) {
	count, err := collection.Count()
	return count > 0, err
}

func TestMongoPersister_Persist_NoPersistName(t *testing.T) {
	if !hasConnection(t) {
		return
	}

	req := definition.Request{}
	res := definition.Response{}

	actions := make(map[string]string)
	actions["write"] = "{{ request.body }}"
	per := definition.Persist{Actions: actions}

	varsProcessor, persister := getMongoProcessor()
	defer func() {
		dropDatabase(persister.Repository) // cleanup database
	}()
	dropDatabase(persister.Repository) // make sure we are working on a clean database

	mock := definition.Mock{Request: req, Response: res, Persist: per}
	varsProcessor.Eval(&req, &mock)

	hasCollections, err := hasCollections(persister.Repository)
	if err != nil {
		t.Error(err)
	}

	if hasCollections {
		t.Error("No collections should be present")
	}
}

func TestMongoPersister_Persist_WithBodyToSave(t *testing.T) {
	if !hasConnection(t) {
		return
	}
	req := definition.Request{}
	res := definition.Response{Body: "{{ persist.entity.content }}", StatusCode: 200}

	actions := make(map[string]string)
	actions["write"] = "{{ request.body }}"
	per := definition.Persist{Actions: actions}

	req.Body = "{ \"test\": 1}"
	per.Entity = "test/testing-1"

	varsProcessor, persister := getMongoProcessor()
	defer func() {
		dropDatabase(persister.Repository) // cleanup database
	}()
	dropDatabase(persister.Repository) // make sure we are working on a clean database

	mock := definition.Mock{Request: req, Response: res, Persist: per}
	varsProcessor.Eval(&req, &mock)

	hasItems, _ := hasCollectionsItems(persister.Repository, "test")

	if !hasItems {
		t.Error("One item should be created")
	} else {

		bodyContent, err := persister.Repository.GetItem("test", "testing-1")
		if err != nil {
			t.Error(err)
		}

		if equal, err := utils.JSONSStringsAreEqual(bodyContent, mock.Response.Body); !equal || err != nil {
			t.Error("File content should match result body", bodyContent, mock.Response.Body)
		}
	}
}

func TestMongoPersister_Persist_WithBodyToSave_WrongName(t *testing.T) {
	if !hasConnection(t) {
		return
	}
	req := definition.Request{}
	res := definition.Response{Body: "{{persist.entity.content}}", StatusCode: 200}

	actions := make(map[string]string)
	actions["write"] = "{{ request.body }}"
	per := definition.Persist{Actions: actions}

	req.Body = "{ \"test\": 1}"
	per.Entity = "test-testing-1"

	varsProcessor, persister := getMongoProcessor()
	defer func() {
		dropDatabase(persister.Repository) // cleanup database
	}()
	dropDatabase(persister.Repository) // make sure we are working on a clean database

	mock := definition.Mock{Request: req, Response: res, Persist: per}
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != "" {
		t.Error("The body should be empty as we have wrong name format. Current body content is: ", res.Body)
	}
}

func TestMongoPersister_Persist_WithNonJSONBodyToSave(t *testing.T) {
	if !hasConnection(t) {
		return
	}
	req := definition.Request{}
	res := definition.Response{Body: "{{persist.entity.content}}", StatusCode: 200}

	actions := make(map[string]string)
	actions["write"] = "{{ request.body }}"
	per := definition.Persist{Actions: actions}

	req.Body = "Body to save"
	per.Entity = "test/testing-1"

	varsProcessor, persister := getMongoProcessor()
	defer func() {
		dropDatabase(persister.Repository) // cleanup database
	}()
	dropDatabase(persister.Repository) // make sure we are working on a clean database

	mock := definition.Mock{Request: req, Response: res, Persist: per}
	varsProcessor.Eval(&req, &mock)

	hasItems, _ := hasCollectionsItems(persister.Repository, "test")

	if !hasItems {
		t.Error("One item should be created")
	} else {

		bodyContent, err := persister.Repository.GetItem("test", "testing-1")
		if err != nil {
			t.Error(err)
		}

		if equal, err := utils.JSONSStringsAreEqual(bodyContent, mock.Response.Body); !equal || err != nil {
			t.Error("File content should match result body", bodyContent, mock.Response.Body)
		}
	}
}

func TestMongoPersister_LoadBody(t *testing.T) {
	if !hasConnection(t) {
		return
	}
	req := definition.Request{}
	res := definition.Response{Body: "{{persist.entity.content}}"}
	per := definition.Persist{Entity: "test/item-1"}

	varsProcessor, persister := getMongoProcessor()
	defer func() {
		dropDatabase(persister.Repository) // cleanup database
	}()
	dropDatabase(persister.Repository) // make sure we are working on a clean database

	mock := definition.Mock{Request: req, Response: res, Persist: per}

	content := "{\"test\": \"body to expect\"}"
	id := "item-1"

	err := persister.Repository.UpsertItem("test", id, content)
	if err != nil {
		t.Error("Item should be saved", err)
	} else {
		varsProcessor.Eval(&req, &mock)

		if equal, err := utils.JSONSStringsAreEqual(mock.Response.Body, content); !equal || err != nil {
			t.Error("Result body and file content should be the same", mock.Response.Body, content)
		}
	}
}

func TestMongoPersister_LoadBody_WrongName(t *testing.T) {
	if !hasConnection(t) {
		return
	}
	req := definition.Request{}
	res := definition.Response{Body: "{{persist.entity.content}}"}
	per := definition.Persist{Entity: "test-item-1"}

	varsProcessor, persister := getMongoProcessor()
	defer func() {
		dropDatabase(persister.Repository) // cleanup database
	}()
	dropDatabase(persister.Repository) // make sure we are working on a clean database

	mock := definition.Mock{Request: req, Response: res, Persist: per}
	varsProcessor.Eval(&req, &mock)

	if res.Body != "{{persist.entity.content}}" {
		t.Error("The body should stay {{persist.entity.content}} as the format of the entity to load is wrong. Current body content is: ", res.Body)
	}
}

func TestMongoPersister_LoadBodyNonJSON(t *testing.T) {
	if !hasConnection(t) {
		return
	}
	req := definition.Request{}
	res := definition.Response{Body: "{{persist.entity.content}}"}
	per := definition.Persist{Entity: "test/item-1"}

	varsProcessor, persister := getMongoProcessor()
	defer func() {
		dropDatabase(persister.Repository) // cleanup database
	}()
	dropDatabase(persister.Repository) // make sure we are working on a clean database

	mock := definition.Mock{Request: req, Response: res, Persist: per}
	content := "Non JSON body to expect"
	id := "item-1"

	err := persister.Repository.UpsertItem("test", id, content)
	if err != nil {
		t.Error("Item should be saved", err)
	} else {
		varsProcessor.Eval(&req, &mock)

		if equal, err := utils.JSONSStringsAreEqual(mock.Response.Body, content); !equal || err != nil {
			t.Error("Result body and file content should be the same", mock.Response.Body, content)
		}
	}
}

func TestMongoPersister_LoadBody_WithAppend(t *testing.T) {
	if !hasConnection(t) {
		return
	}
	req := definition.Request{}
	res := definition.Response{Body: "{{persist.entity.content}}"}

	appendText := "{\"append\":1}"

	actions := make(map[string]string)
	actions["append"] = appendText
	per := definition.Persist{Entity: "test/item-1", Actions: actions}

	varsProcessor, persister := getMongoProcessor()
	defer func() {
		dropDatabase(persister.Repository) // cleanup database
	}()
	dropDatabase(persister.Repository) // make sure we are working on a clean database

	mock := definition.Mock{Request: req, Response: res, Persist: per}

	content := "{\"test\": \"body to expect\"}"
	id := "item-1"

	err := persister.Repository.UpsertItem("test", id, content)
	if err != nil {
		t.Error("Item should be saved", err)
	} else {
		varsProcessor.Eval(&req, &mock)

		if equal, err := utils.JSONSStringsAreEqual(mock.Response.Body, utils.JoinJSON(content, appendText)); !equal || err != nil {
			t.Error("Result body and file content plus bodyAppend should be the same", mock.Response.Body, content, appendText)
		}
	}
}

func TestMongoPersister_LoadBody_NotFound(t *testing.T) {
	if !hasConnection(t) {
		return
	}
	req := definition.Request{}
	res := definition.Response{Body: "{{persist.entity.content}}"}
	per := definition.Persist{Entity: "test/item-1"}

	varsProcessor, persister := getMongoProcessor()
	defer func() {
		dropDatabase(persister.Repository) // cleanup database
	}()
	dropDatabase(persister.Repository) // make sure we are working on a clean database

	mock := definition.Mock{Request: req, Response: res, Persist: per}
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != "" {
		t.Error("Result body should be empty", mock.Response.Body)
	} else if mock.Response.StatusCode != 404 {
		t.Error("Status code should be 404", mock.Response.StatusCode)
	}
}

func TestMongoPersister_Sequence(t *testing.T) {
	if !hasConnection(t) {
		return
	}

	persister := persist.NewMongoPersister(mongoTestURL)
	defer func() {
		dropDatabase(persister.Repository) // cleanup database
	}()
	dropDatabase(persister.Repository) // make sure we are working on a clean database

	result, err := persister.GetSequence("test", 0)
	if err != nil {
		t.Error(err)
	}
	if result != 0 {
		t.Error("The result should be 0 as there's no record for the test sequence", result)
	}

	result, err = persister.GetSequence("test", 1)
	if err != nil {
		t.Error(err)
	}
	if result != 1 {
		t.Error("The result should be 1 as the previous value should have been 0", result)
	}

	result, err = persister.GetSequence("test", 4)
	if err != nil {
		t.Error(err)
	}
	if result != 5 {
		t.Error("The result should be 5 as the previous value should have been 1", result)
	}
}

func TestMongoPersister_SetValue_GetValue(t *testing.T) {
	if !hasConnection(t) {
		return
	}

	persister := persist.NewMongoPersister(mongoTestURL)
	defer func() {
		dropDatabase(persister.Repository) // cleanup database
	}()
	dropDatabase(persister.Repository) // make sure we are working on a clean database

	// not existing get
	result, err := persister.GetValue("test")
	if err == nil {
		t.Error("We should have error as the value is not existing")
	}

	// not existing set
	err = persister.SetValue("test", "123")
	if err != nil {
		t.Error(err)
	}

	result, err = persister.GetValue("test")
	if err != nil {
		t.Error(err)
	}
	if result != "123" {
		t.Error("The result should be 0 as there's no record for the test sequence", result)
	}

	// existing set
	err = persister.SetValue("test", "456")
	if err != nil {
		t.Error(err)
	}

	result, err = persister.GetValue("test")
	if err != nil {
		t.Error(err)
	}
	if result != "456" {
		t.Error("The result should be 0 as there's no record for the test sequence", result)
	}
}
