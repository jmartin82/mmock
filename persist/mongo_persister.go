package persist

import (
	"errors"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/jmartin82/mmock/utils"
	"github.com/ryanuber/go-glob"
)

var (
	errWrongNameFormat = errors.New("The name of the persist item should be in the following format {collectionName}/{itemId}")
)

//MongoPersister persists body in mongo
type MongoPersister struct {
	Repository MongoRepository
}

func (mp MongoPersister) GetName() string {
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

func (mp MongoPersister) ReadCollection(name string) (string, error) {
	log.Printf("Reading collection: %s\n", name)
	itemsInCollection := mp.getCollectionItems(name)

	contents := []string{}
	allJSON := true

	var keys []string
	for k := range itemsInCollection {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		content := itemsInCollection[key]
		if allJSON {
			allJSON = utils.IsJSON(content)
		}
		contents = append(contents, content)
	}
	if allJSON {
		return "[" + strings.Join(contents, ",") + "]", nil
	}
	return strings.Join(contents, "\n"), nil
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

func (mp MongoPersister) DeleteCollection(name string) error {
	log.Printf("Deleting collection: %s\n", name)
	itemsInCollection := mp.getCollectionItems(name)

	for key, _ := range itemsInCollection {
		collectionName, id, err := mp.getItemInfo(key)
		if err == nil {
			log.Printf("Deleting entity %s from collection %s\n", id, collectionName)
			err = mp.Repository.DeleteItem(collectionName, id)
		}
	}
	return nil
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

func (mp MongoPersister) getItemsInCollection(collectionName string) map[string]string {
	collection := collectionName
	if strings.Index(collection, "/") == 0 {
		collection = collection[1:]
	}
	filter := ""
	if i := strings.Index(collection, "/"); i > 0 {
		filter = collection[(i + 1):]
		collection = collection[0:i]
	}
	items := map[string]string{}

	if collection == "" {
		return items
	}

	regex, regexError := regexp.Compile(filter)

	dataItems, err := mp.Repository.GetAllItems(collection)

	if err != nil {
		log.Println(err)
		return items
	}

	for key, value := range dataItems {
		if (filter == "") || glob.Glob(filter, key) || (regexError == nil && regex.MatchString(key)) {
			items[collection+"/"+key] = value
		}
	}

	return items
}

func (mp MongoPersister) getItemsList(name string) map[string]string {
	if strings.Index(name, ",") == 0 {
		name = name[1:] // remove the starting comma
	}
	keys := strings.Split(name, ",")
	items := map[string]string{}
	for _, key := range keys {
		collection, id, err := mp.getItemInfo(key)
		if err == nil {
			value, err := mp.Repository.GetItem(collection, id)
			if err == nil {
				items[key] = value
			}
		}
	}
	return items
}

func (mp MongoPersister) getCollectionItems(name string) map[string]string {
	if strings.Index(name, ",") > -1 {
		return mp.getItemsList(name)
	}

	return mp.getItemsInCollection(name)
}

//NewMongoPersister creates a new MongoPersister
func NewMongoPersister(connectionString string) *MongoPersister {

	mongoRepo := NewMongoRepository(connectionString)
	result := MongoPersister{Repository: *mongoRepo}

	return &result
}
