package persist

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/parse"
)

//FileBodyPersister persists body in file
type FileBodyPersister struct {
	PersistPath string
	Parser      parse.ResponseParser
}

//Persist the body of the response to fiel if needed
func (fbp FileBodyPersister) Persist(per *definition.Persist, req *definition.Request, res *definition.Response) {
	if per.Name == "" {
		return
	}

	per.Name = fbp.Parser.ReplaceVars(req, per.Name)

	filePath := path.Join(fbp.PersistPath, per.Name)
	fileDir := path.Dir(filePath)

	if per.Delete {
		os.Remove(filePath)
	} else {
		fileContent := []byte(res.Body)
		err := os.MkdirAll(fileDir, 0644)
		if fbp.checkForFileWriteError(err, res) == nil {
			err = ioutil.WriteFile(filePath, fileContent, 0644)
			fbp.checkForFileWriteError(err, res)
		}
	}
}

//LoadBody loads the response body from the persisted file
func (fbp FileBodyPersister) LoadBody(req *definition.Request, res *definition.Response) {
	a := "test"
	_ = a
}

func (fbp FileBodyPersister) checkForFileWriteError(err error, res *definition.Response) error {
	if err != nil {
		log.Print(err)
		res.Body = err.Error()
		res.StatusCode = 500
	}
	return err
}

//NewFileBodyPersister creates a new FileBodyPersister
func NewFileBodyPersister(persistPath string, parser parse.ResponseParser) *FileBodyPersister {
	result := FileBodyPersister{PersistPath: persistPath, Parser: parser}

	err := os.MkdirAll(result.PersistPath, 0644)
	if err != nil {
		panic(err)
	}

	return &result
}
