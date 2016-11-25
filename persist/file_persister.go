package persist

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/jmartin82/mmock/utils"
	"github.com/ryanuber/go-glob"
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

func (fp FilePersister) ReadCollection(name string) (string, error) {
	log.Printf("Reading collection: %s\n", name)
	filesInCollection := fp.getCollectionFiles(name)

	contents := []string{}
	allJSON := true

	sort.Strings(filesInCollection)

	for _, file := range filesInCollection {
		if fileContent, err := ioutil.ReadFile(file); err == nil {
			stringContent := string(fileContent)
			if allJSON {
				allJSON = utils.IsJSON(stringContent)
			}
			contents = append(contents, stringContent)
		}
	}
	if allJSON {
		return "[" + strings.Join(contents, ",") + "]", nil
	}
	return strings.Join(contents, "\n"), nil
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

func (fp FilePersister) DeleteCollection(name string) error {
	log.Printf("Deleting collection: %s\n", name)
	filesInCollection := fp.getCollectionFiles(name)

	for _, file := range filesInCollection {
		err := os.Remove(file)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
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

func (fp FilePersister) getFilesInCollection(collectionName string) []string {
	folder := collectionName
	if strings.Index(folder, "/") == 0 {
		folder = folder[1:]
	}
	filter := ""
	if i := strings.Index(folder, "/"); i > 0 {
		filter = folder[(i + 1):]
		folder = folder[0:i]
	}
	fullFolderPath := fp.getFolderPath(folder)

	filesList := []string{}
	if fullFolderPath == "" {
		return []string{}
	}

	regex, regexError := regexp.Compile(filter)

	filepath.Walk(fullFolderPath, func(filePath string, fileInfo os.FileInfo, err error) error {
		if !fileInfo.IsDir() {
			if len(filePath) > len(fullFolderPath) {
				relativeFilePath := filePath[(len(fullFolderPath) + 1):]

				if (filter == "") || glob.Glob(filter, relativeFilePath) || (regexError == nil && regex.MatchString(relativeFilePath)) {
					filesList = append(filesList, filePath)
				}
			}
		}
		return nil
	})
	return filesList
}

func (fp FilePersister) getFilesList(name string) []string {
	if strings.Index(name, ",") == 0 {
		name = name[1:] // remove the starting comma
	}
	fileNames := strings.Split(name, ",")
	files := []string{}
	for _, fileName := range fileNames {
		pathToFile, _ := fp.getFilePath(fileName)
		if pathToFile != "" {
			files = append(files, pathToFile)
		}
	}
	return files
}

func (fp FilePersister) getCollectionFiles(name string) []string {
	if strings.Index(name, ",") > -1 {
		return fp.getFilesList(name)
	}

	return fp.getFilesInCollection(name)
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
