package file

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/parse"
)

//FileBodyPersister persists body in file
type FileBodyPersister struct {
	PersistPath string
	Parser      parse.ResponseParser
}

//Persist the body of the response to file if needed
func (fbp FileBodyPersister) Persist(per *definition.Persist, req *definition.Request, res *definition.Response) bool {
	result := true

	if per.Name == "" {
		return result
	}

	per.Name = fbp.Parser.ReplaceVars(req, res, per.Name)

	pathToFile, fileDir := fbp.getFilePath(per.Name, req, res)
	if pathToFile == "" {
		return false
	}

	if per.Delete {
		os.Remove(pathToFile)
	} else {
		fileContent := []byte(res.Body)
		err := os.MkdirAll(fileDir, 0755)
		result = (fbp.checkForError(err, res) == nil)
		if result {
			err = ioutil.WriteFile(pathToFile, fileContent, 0755)
			result = (fbp.checkForError(err, res) == nil)
		}
	}

	return result
}

func (fbp FileBodyPersister) getFilePath(fileName string, req *definition.Request, res *definition.Response) (pathToFile string, fileDir string) {
	pathToFile = path.Join(fbp.PersistPath, fileName)
	fileDir = path.Dir(pathToFile)

	var err error
	pathToFile, err = filepath.Abs(pathToFile)
	if fbp.checkForError(err, res) != nil {
		return "", ""
	}
	if !strings.HasPrefix(pathToFile, fbp.PersistPath) {
		errorText := fmt.Sprintf("File path not under the persist path. FilePath: %s, PersistPath %s", pathToFile, fbp.PersistPath)
		fbp.checkForError(errors.New(errorText), res)
		return "", ""
	}

	return pathToFile, fileDir
}

//LoadBody loads the response body from the persisted file
func (fbp FileBodyPersister) LoadBody(req *definition.Request, res *definition.Response) {
	res.Persisted.Name = fbp.Parser.ReplaceVars(req, res, res.Persisted.Name)

	pathToFile, _ := fbp.getFilePath(res.Persisted.Name, req, res)
	if pathToFile == "" {
		return
	}

	_, err := os.Stat(pathToFile)

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
		fileContent, err := ioutil.ReadFile(pathToFile)
		if fbp.checkForError(err, res) == nil {
			res.Body = fbp.Parser.ParseBody(req, res, string(fileContent), res.Persisted.BodyAppend)
		}
	}
}

func (fbp FileBodyPersister) checkForError(err error, res *definition.Response) error {
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

	err := os.MkdirAll(result.PersistPath, 0755)
	if err != nil {
		panic(err)
	}

	return &result
}
