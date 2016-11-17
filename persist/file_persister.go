package persist

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//FilePersister persists body in file
type FilePersister struct {
	PersistPath string
}

func (fbp FilePersister) GetName() string {
	return "file"
}

func (fbp FilePersister) Read(name string) (string, error) {

	pathToFile, _ := fbp.getFilePath(name)
	log.Printf("Reading entity: %s\n", pathToFile)
	if _, err := os.Stat(pathToFile); err != nil {
		log.Printf("Error reading the entity (%s)\n", err)
		return "", err
	}

	content, err := ioutil.ReadFile(pathToFile)
	return string(content), err
}
func (fbp FilePersister) Write(name, content string) error {
	pathToFile, fileDir := fbp.getFilePath(name)
	fileContent := []byte(content)
	err := os.MkdirAll(fileDir, 0755)
	if err == nil {
		err = ioutil.WriteFile(pathToFile, fileContent, 0755)
	}
	return err
}
func (fbp FilePersister) Delete(name string) error {
	pathToFile, _ := fbp.getFilePath(name)
	return os.Remove(pathToFile)
}

func (fbp FilePersister) getFilePath(fileName string) (pathToFile string, fileDir string) {
	pathToFile = path.Join(fbp.PersistPath, fileName)
	fileDir = path.Dir(pathToFile)

	var err error
	pathToFile, err = filepath.Abs(pathToFile)
	if err != nil {
		return "", ""
	}
	if !strings.HasPrefix(pathToFile, fbp.PersistPath) {
		log.Printf("File path not under the persist path. FilePath: %s, PersistPath %s", pathToFile, fbp.PersistPath)
		return "", ""
	}

	return pathToFile, fileDir
}

//NewFilePersister creates a new FilePersister
func NewFilePersister(persistPath string) *FilePersister {
	result := FilePersister{PersistPath: persistPath}

	err := os.MkdirAll(result.PersistPath, 0755)
	if err != nil {
		panic(err)
	}

	return &result
}
