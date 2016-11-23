package persist

import (
	"errors"
	"log"
	"strings"
)

var (
	errWrongNameFormat = errors.New("The name of the persist item should be in the following format {collectionName}/{itemId}")
)

//MongoPersister persists body in mongo
type MongoPersister struct {
	Repository MongoRepository
}

func (fbp MongoPersister) GetName() string {
	return "mongo"
}

func (mp MongoPersister) Read(name string) (string, error) {

	collectionName, id, err := mp.getItemInfo(name)
	if err != nil {
		return "", err
	}
	log.Printf("Reading entity %s from collection %s\n", id, collectionName)

	content, err := mp.Repository.GetItem(collectionName, id)
	return content, err
}

func (mp MongoPersister) Write(name, content string) error {
	collectionName, id, err := mp.getItemInfo(name)
	if err == nil {
		log.Printf("Writing entity %s from collection %s\n", id, collectionName)
		err = mp.Repository.UpsertItem(collectionName, id, content)
	}
	return err
}

func (mp MongoPersister) Delete(name string) error {
	collectionName, id, err := mp.getItemInfo(name)
	if err == nil {
		log.Printf("Deleting entity %s from collection %s\n", id, collectionName)
		err = mp.Repository.DeleteItem(collectionName, id)
	}
	return err
}

func (mp MongoPersister) getItemInfo(name string) (collectionName string, id string, err error) {
	// remove starting slash
	if i := strings.Index(name, "/"); i == 0 {
		name = name[1:]
	}
	if i := strings.Index(name, "/"); i > 0 {
		return name[0:i], name[(i + 1):], nil
	} else {
		return "", "", errWrongNameFormat
	}
}

//NewMongoPersister creates a new MongoPersister
func NewMongoPersister(connectionString string) *MongoPersister {

	mongoRepo := NewMongoRepository(connectionString)
	result := MongoPersister{Repository: *mongoRepo}

	return &result
}
