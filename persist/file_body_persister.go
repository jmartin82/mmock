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
func (fbp FileBodyPersister) Persist(per *definition.Persist, req *definition.Request, res *definition.Response) bool {
	result := true

	if per.Name == "" {
		return result
	}

	per.Name = fbp.Parser.ReplaceVars(req, res, per.Name)

	filePath := path.Join(fbp.PersistPath, per.Name)
	fileDir := path.Dir(filePath)

	if per.Delete {
		os.Remove(filePath)
	} else {
		fileContent := []byte(res.Body)
		err := os.MkdirAll(fileDir, 0644)
		result = (checkForError(err, res) == nil)
		if result {
			err = ioutil.WriteFile(filePath, fileContent, 0644)
			result = (checkForError(err, res) == nil)
		}
	}

	return result
}

func (fbp FileBodyPersister) getFilePath(fileName string, req *definition.Request, res *definition.Response) (filePath string, fileDir string) {
	fileName = fbp.Parser.ReplaceVars(req, res, fileName)

	filePath = path.Join(fbp.PersistPath, fileName)
	fileDir = path.Dir(filePath)

	return filePath, fileDir
}

//LoadBody loads the response body from the persisted file
func (fbp FileBodyPersister) LoadBody(req *definition.Request, res *definition.Response) {
	filePath, _ := fbp.getFilePath(res.Persisted.Name, req, res)

	var _, err = os.Stat(filePath)

	// use notFound info
	if os.IsNotExist(err) {
		res.Body = "Not Found"
		if res.Persisted.NotFound.Body != "" {
			res.Body = fbp.Parser.ParseBody(req, res, res.Persisted.NotFound.Body, res.Persisted.NotFound.BodyAppend)
		}
		res.StatusCode = 404
		if res.Persisted.NotFound.StatusCode != 0 {
			res.StatusCode = res.Persisted.NotFound.StatusCode
		}
	} else {
		fileContent, err := ioutil.ReadFile(filePath)
		if checkForError(err, res) == nil {
			res.Body = fbp.Parser.ParseBody(req, res, string(fileContent), res.Persisted.BodyAppend)
		}
	}
}

func checkForError(err error, res *definition.Response) error {
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
