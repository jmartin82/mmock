package persist

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/jmartin82/mmock/parse"
	"github.com/tidwall/sjson"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var errMongoDatabaseNotPassed = errors.New("Mongo database not passed. Please add the database at the end of the connection string e.g. /DatabasName ")

//MongoRepository saves and loads items from mongo
type MongoRepository struct {
	ConnectionInfo mgo.DialInfo
}

func (mr MongoRepository) ConnectMongo() (session *mgo.Session, err error) {
	session, err = mgo.DialWithInfo(&mr.ConnectionInfo)
	if err == nil {
		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)
	}
	return session, err
}

//GetItem gets the result string rom interface
func (mr MongoRepository) GetItem(collectionName string, id string) (string, error) {

	session, err := mr.ConnectMongo()
	if err != nil {
		return "", err
	}

	defer session.Close()

	collection := session.DB(mr.ConnectionInfo.Database).C(collectionName)

	var result interface{}
	err = collection.FindId(id).One(&result)
	if err != nil {
		return "", err
	}

	byteResult, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	resultString := string(byteResult)

	// remove _id property from the result json string
	resultString, err = sjson.Delete(resultString, "_id")
	if err != nil {
		return "", err
	}

	return parse.UnWrapNonJSONStringIfNeeded(resultString), nil
}

//DeleteItem deletes an item from a collection
func (mr MongoRepository) DeleteItem(collectionName string, id string) error {

	session, err := mr.ConnectMongo()
	if err != nil {
		return err
	}

	defer session.Close()

	collection := session.DB(mr.ConnectionInfo.Database).C(collectionName)

	return collection.RemoveId(id)
}

//UpsertItem inserts or updates item with a given id in a given collection
func (mr MongoRepository) UpsertItem(collectionName string, id string, body string) error {
	session, err := mr.ConnectMongo()
	if err != nil {
		return err
	}

	defer session.Close()

	collection := session.DB(mr.ConnectionInfo.Database).C(collectionName)

	body, err = parse.WrapNonJSONStringIfNeeded(body)
	if err != nil {
		return err
	}

	var parsedData interface{}
	json.Unmarshal([]byte(body), &parsedData)

	upsertdata := bson.M{"$set": parsedData}

	_, err = collection.UpsertId(id, upsertdata)
	return err
}

//NewMongoRepository creates a new MongoRepository
func NewMongoRepository(connectionString string) *MongoRepository {

	dialInfo, err := mgo.ParseURL(connectionString)
	if err != nil {
		log.Print(err)
	}

	if dialInfo.Database == "" {
		panic(errMongoDatabaseNotPassed)
	}

	if dialInfo.Timeout == 0 {
		dialInfo.Timeout = 10 * time.Second
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	result := MongoRepository{ConnectionInfo: *dialInfo}

	return &result
}
