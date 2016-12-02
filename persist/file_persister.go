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

func (fp FilePersister) GetName() string {
	return "file"
}

func (fp FilePersister) Read(name string) (string, error) {

	pathToFile, _ := fp.getFilePath(name)
	log.Printf("Reading entity: %s\n", pathToFile)
	if _, err := os.Stat(pathToFile); err != nil {
		log.Printf("Error reading the entity (%s)\n", err)
		return "", err
	}

	content, err := ioutil.ReadFile(pathToFile)
	return string(content), err
}

func (fp FilePersister) Write(name, content string) error {
	pathToFile, fileDir := fp.getFilePath(name)
	fileContent := []byte(content)
	err := os.MkdirAll(fileDir, 0755)
	if err == nil {
		err = ioutil.WriteFile(pathToFile, fileContent, 0755)
	}
	return err
}

func (fp FilePersister) Delete(name string) error {
	pathToFile, _ := fp.getFilePath(name)
	return os.Remove(pathToFile)
}

func (fp FilePersister) getFilePath(fileName string) (pathToFile string, fileDir string) {
	pathToFile = path.Join(fp.PersistPath, fileName)
	fileDir = path.Dir(pathToFile)

	var err error
	pathToFile, err = filepath.Abs(pathToFile)
	if err != nil {
		return "", ""
	}
	if !strings.HasPrefix(pathToFile, fp.PersistPath) {
		log.Printf("File path not under the persist path. FilePath: %s, PersistPath %s", pathToFile, fp.PersistPath)
		return "", ""
	}

	return pathToFile, fileDir
}

func (fp FilePersister) getFolderPath(folderName string) (pathToFolder string) {
	pathToFolder = path.Join(fp.PersistPath, folderName)

	var err error
	pathToFolder, err = filepath.Abs(pathToFolder)
	if err != nil {
		return ""
	}
	if !strings.HasPrefix(pathToFolder, fp.PersistPath) {
		log.Printf("Folder path not under the persist path. FilePath: %s, PersistPath %s", pathToFolder, fp.PersistPath)
		return ""
	}

	return pathToFolder
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
