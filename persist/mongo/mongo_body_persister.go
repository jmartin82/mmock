package mongo

// import (
// 	"encoding/json"
// 	"errors"
// 	"log"
// 	"strings"

// 	"github.com/jmartin82/mmock/definition"
// 	"github.com/jmartin82/mmock/parse"
// 	"github.com/tidwall/sjson"
// )

// var (
// 	errWrongNameFormat = errors.New("The name of the persist item should be in the following format {collectionName}/{itemId}")
// )

// //MongoBodyPersister persists body in file
// type MongoBodyPersister struct {
// 	Parser     parse.ResponseParser
// 	Repository MongoRepository
// }

// //Persist the body of the response to mongo if needed
// func (mbp MongoBodyPersister) Persist(per *definition.Persist, req *definition.Request, res *definition.Response) bool {
// 	result := true

// 	if per.Name == "" {
// 		return result
// 	}

// 	per.Name = mbp.Parser.ReplaceVars(req, res, per.Name)

// 	collectionName, id := mbp.getItemInfo(per.Name, res)

// 	if collectionName != "" {
// 		var err error
// 		if per.Delete {
// 			err = mbp.Repository.DeleteItem(collectionName, id)
// 		} else {
// 			err = mbp.Repository.UpsertItem(collectionName, id, res.Body)
// 		}

// 		if err != nil && err.Error() == "not found" {
// 			err = nil
// 		}
// 		result = mbp.checkForError(err, res) != nil
// 	} else {
// 		result = false
// 	}

// 	return result
// }

// func (mbp MongoBodyPersister) getItemInfo(name string, res *definition.Response) (collectionName string, id string) {
// 	if i := strings.Index(name, "/"); i > 0 {
// 		collectionName = name[0:i]
// 		id = name[(i + 1):]
// 	} else {
// 		mbp.checkForError(errWrongNameFormat, res)
// 	}

// 	return collectionName, id
// }

// //GetResultString gets the result string rom interface
// func (mbp MongoBodyPersister) GetResultString(result interface{}) (string, error) {
// 	byteResult, err := json.Marshal(result)
// 	if err != nil {
// 		return "", err
// 	}

// 	resultString := string(byteResult)

// 	// remove _id property from the result json string
// 	resultString, _ = sjson.Delete(resultString, "_id")

// 	return resultString, nil
// }

// //LoadBody loads the response body from the corresponding mongo collection
// func (mbp MongoBodyPersister) LoadBody(req *definition.Request, res *definition.Response) {
// 	res.Persisted.Name = mbp.Parser.ReplaceVars(req, res, res.Persisted.Name)

// 	collectionName, id := mbp.getItemInfo(res.Persisted.Name, res)
// 	if collectionName == "" {
// 		return
// 	}

// 	resultString, err := mbp.Repository.GetItem(collectionName, id)

// 	if err != nil {
// 		errorString := err.Error()
// 		if errorString != "not found" {
// 			mbp.checkForError(err, res)
// 			return
// 		} else if err != nil && err.Error() == "not found" {
// 			res.Body = "Not Found"
// 			if res.Persisted.NotFound.Body != "" {
// 				res.Body = mbp.Parser.ParseBody(req, res, res.Persisted.NotFound.Body, res.Persisted.NotFound.BodyAppend)
// 			}
// 			res.StatusCode = 404
// 			if res.Persisted.NotFound.StatusCode != 0 {
// 				res.StatusCode = res.Persisted.NotFound.StatusCode
// 			}
// 		}
// 	} else {
// 		res.Body = mbp.Parser.ParseBody(req, res, resultString, res.Persisted.BodyAppend)
// 	}
// }

// func (mbp MongoBodyPersister) checkForError(err error, res *definition.Response) error {
// 	if err != nil {
// 		log.Print(err)
// 		res.Body = err.Error()
// 		res.StatusCode = 500
// 	}
// 	return err
// }

// //NewMongoBodyPersister creates a new MongoBodyPersister
// func NewMongoBodyPersister(connectionString string, parser parse.ResponseParser) *MongoBodyPersister {

// 	mongoRepo := NewMongoRepository(connectionString)
// 	result := MongoBodyPersister{Parser: parser, Repository: *mongoRepo}

// 	return &result
// }
