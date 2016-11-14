package persist

import (
	"encoding/json"
	"log"
	"strings"

	"errors"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/parse"
	"github.com/tidwall/sjson"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//MongoBodyPersister persists body in file
type MongoBodyPersister struct {
	ConnectionInfo mgo.DialInfo
	Parser         parse.ResponseParser
}

//Persist the body of the response to mongo if needed
func (mbp MongoBodyPersister) Persist(per *definition.Persist, req *definition.Request, res *definition.Response) bool {
	result := true

	if per.Name == "" {
		return result
	}

	per.Name = mbp.Parser.ReplaceVars(req, res, per.Name)

	collectionName, id := mbp.getItemInfo(per.Name, res)

	if collectionName != "" {
		var err error
		if per.Delete {
			err = mbp.deleteItem(collectionName, id)
		} else {
			err = mbp.saveItem(collectionName, id, res.Body)
		}

		if err != nil && err.Error() == "not found" {
			err = nil
		}
		result = mbp.checkForError(err, res) != nil
	} else {
		result = false
	}

	return result
}

func (mbp MongoBodyPersister) getItemInfo(name string, res *definition.Response) (collectionName string, id string) {
	if i := strings.Index(name, "/"); i > 0 {
		collectionName = name[0:i]
		id = name[(i + 1):]
	} else {
		mbp.checkForError(errors.New("The name of the persist item should be in the following format {collectionName}/{collectionId}"), res)
	}

	return collectionName, id
}

func (mbp MongoBodyPersister) connectMongo() (session *mgo.Session, err error) {
	session, err = mgo.DialWithInfo(&mbp.ConnectionInfo)
	if err == nil {
		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)
	}
	return session, err
}

func (mbp MongoBodyPersister) saveItem(collectionName string, id string, body string) error {
	session, err := mbp.connectMongo()
	if err != nil {
		return err
	}

	defer session.Close()

	collection := session.DB(mbp.ConnectionInfo.Database).C(collectionName)

	var parsedData interface{}
	json.Unmarshal([]byte(body), &parsedData)

	upsertdata := bson.M{"$set": parsedData}

	_, err = collection.UpsertId(id, upsertdata)
	return err
}

func (mbp MongoBodyPersister) deleteItem(collectionName string, id string) error {
	session, err := mbp.connectMongo()
	if err != nil {
		return err
	}

	defer session.Close()

	collection := session.DB(mbp.ConnectionInfo.Database).C(collectionName)

	return collection.RemoveId(id)
}

func (mbp MongoBodyPersister) getResultString(result interface{}) (string, error) {
	byteResult, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	resultString := string(byteResult)

	// remove _id property from the result json string
	resultString, _ = sjson.Delete(resultString, "_id")

	return resultString, nil
}

//LoadBody loads the response body from the corresponding mongo collection
func (mbp MongoBodyPersister) LoadBody(req *definition.Request, res *definition.Response) {
	res.Persisted.Name = mbp.Parser.ReplaceVars(req, res, res.Persisted.Name)

	collectionName, id := mbp.getItemInfo(res.Persisted.Name, res)
	if collectionName == "" {
		return
	}

	session, err := mbp.connectMongo()
	if mbp.checkForError(err, res) != nil {
		return
	}

	defer session.Close()

	collection := session.DB(mbp.ConnectionInfo.Database).C(collectionName)

	var result interface{}
	err = collection.FindId(id).One(&result)

	if err != nil {
		errorString := err.Error()
		if errorString != "not found" {
			mbp.checkForError(err, res)
			return
		} else if err != nil && err.Error() == "not found" {
			res.Body = "Not Found"
			if res.Persisted.NotFound.Body != "" {
				res.Body = mbp.Parser.ParseBody(req, res, res.Persisted.NotFound.Body, res.Persisted.NotFound.BodyAppend)
			}
			res.StatusCode = 404
			if res.Persisted.NotFound.StatusCode != 0 {
				res.StatusCode = res.Persisted.NotFound.StatusCode
			}
		}
	} else {
		resultString, err := mbp.getResultString(result)
		if mbp.checkForError(err, res) == nil {
			res.Body = mbp.Parser.ParseBody(req, res, resultString, res.Persisted.BodyAppend)
		}
	}
}

func (mbp MongoBodyPersister) checkForError(err error, res *definition.Response) error {
	if err != nil {
		log.Print(err)
		res.Body = err.Error()
		res.StatusCode = 500
	}
	return err
}

//NewMongoBodyPersister creates a new MongoBodyPersister
func NewMongoBodyPersister(mongoConnectionString string, parser parse.ResponseParser) *MongoBodyPersister {

	dialInfo, err := mgo.ParseURL(mongoConnectionString)
	if err != nil {
		log.Print(err)
	}

	if dialInfo.Database == "" {
		panic(errors.New("Mongo database not passed. Please add the database at the end of the connection string e.g. /DatabasName "))
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	result := MongoBodyPersister{ConnectionInfo: *dialInfo, Parser: parser}

	return &result
}
