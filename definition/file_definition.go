package definition

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/fsnotify/fsnotify"
	"github.com/jmartin82/mmock/logging"
)

//ErrNotFoundPath error from missing or configuration path
var ErrNotFoundPath = errors.New("Configuration path not found")

//NewFileDefinition file definition constructor
func NewFileDefinition(path string, updatesCh chan []Mock) *FileDefinition {
	return &FileDefinition{
		Path:          path,
		Updates:       updatesCh,
		ConfigReaders: []ConfigReader{},
	}
}

//FileDefinition this struct contains the path of definition and some config readers
type FileDefinition struct {
	Path          string
	Updates       chan []Mock
	ConfigReaders []ConfigReader
}

//PrioritySort mock array sorted by priority
type PrioritySort []Mock

func (s PrioritySort) Len() int {
	return len(s)
}
func (s PrioritySort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s PrioritySort) Less(i, j int) bool {
	return s[i].Control.Priority > s[j].Control.Priority
}

func (fd *FileDefinition) existsConfigPath(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (fd *FileDefinition) getConfigFiles(path string) []string {
	filesList := []string{}
	filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
		if !fileInfo.IsDir() {
			filesList = append(filesList, filePath)
		}
		return nil
	})

	return filesList
}

func (fd *FileDefinition) readMock(filename string) (Mock, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return Mock{}, err
	}
	m := Mock{}
	err = json.Unmarshal(buf, &m)
	if err != nil {
		logging.Printf("Invalid mock definition in: %s\n", filename)
		return Mock{}, err
	}
	m.Name = filepath.Base(filename)
	return m, nil
}

//WatchDir start the watching process to detect any change on defintions
func (fd *FileDefinition) WatchDir() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logging.Println("Hot mock file changing not available.")
		return
	}

	err = watcher.Add(fd.Path)
	if err != nil {
		logging.Printf("Hot mock file changing not available in folder: %s\n", fd.Path)
		return
	}

	if err = filepath.Walk(fd.Path, func(path string, fileInfo os.FileInfo, err error) error {
		if fileInfo.IsDir() {
			err = watcher.Add(path)
			if err != nil {
				logging.Printf("Hot mock file changing not available in folder: %s\n", fd.Path)
				return err
			}
		}
		return nil
	}); err != nil {
		return
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Remove == fsnotify.Remove {

					logging.Println("Changes detected in mock definitions")
					fd.Updates <- fd.ReadMocksDefinition()

				}
			case err := <-watcher.Errors:
				logging.Fatal(err)
			}
		}
	}()
}

//AddConfigReader allows append new readers to able load different config files
func (fd *FileDefinition) AddConfigReader(reader ConfigReader) {
	fd.ConfigReaders = append(fd.ConfigReaders, reader)
}

//ReadMocksDefinition reads all definitions and return an array of valid mocks
func (fd *FileDefinition) ReadMocksDefinition() []Mock {

	if !fd.existsConfigPath(fd.Path) {
		logging.Fatalf(ErrNotFoundPath.Error())
	}

	mocks := []Mock{}
	for _, file := range fd.getConfigFiles(fd.Path) {
		for _, reader := range fd.ConfigReaders {
			if reader.CanRead(file) {
				if mockDef, err := reader.Read(file); err == nil {
					mockDef.Name = filepath.Base(file)
					mocks = append(mocks, mockDef)

				}
				break
			}

		}

	}

	sort.Sort(PrioritySort(mocks))

	return mocks
}
