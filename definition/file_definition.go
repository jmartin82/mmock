package definition

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/fsnotify/fsnotify"
)

//ErrNotFoundPath error from missing or configuration path
var ErrNotFoundPath = errors.New("Configuration path not found")

//ErrNotValidConfigReaderFound we don't have any config reader valid for this file
var ErrNotValidConfigReaderFound = errors.New("Not valid config reader found")

//ErrInvalidMockDefinition the file contains an invalid mock definition
var ErrInvalidMockDefinition = errors.New("Invalid mock definition")

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
	watcher       *fsnotify.Watcher
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

func (fd *FileDefinition) UnWatchDir() {
	if fd.watcher != nil {
		fd.watcher.Close()
	}
}

//WatchDir start the watching process to detect any change on defintions
func (fd *FileDefinition) WatchDir() {
	var err error
	fd.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Println("Hot mock file changing not available.")
		return
	}

	err = fd.watcher.Add(fd.Path)
	if err != nil {
		log.Printf("Hot mock file changing not available in folder: %s\n", fd.Path)
		return
	}

	if err = filepath.Walk(fd.Path, func(path string, fileInfo os.FileInfo, err error) error {
		if fileInfo.IsDir() {
			err = fd.watcher.Add(path)
			if err != nil {
				log.Printf("Hot mock file changing not available in folder: %s\n", fd.Path)
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
			case event := <-fd.watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Remove == fsnotify.Remove {

					log.Println("Changes detected in mock definitions")
					fd.Updates <- fd.ReadMocksDefinition()

				}
			}
		}
	}()
}

//AddConfigReader allows append new readers to able load different config files
func (fd *FileDefinition) AddConfigReader(reader ConfigReader) {
	fd.ConfigReaders = append(fd.ConfigReaders, reader)
}

func (fd *FileDefinition) getMockFromFile(filename string) (Mock, error) {
	for _, reader := range fd.ConfigReaders {
		if reader.CanRead(filename) {
			buf, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Printf("Invalid mock definition in: %s\n", filename)
				return Mock{}, ErrInvalidMockDefinition
			}
			return reader.Read(buf)
		}

	}
	return Mock{}, ErrNotValidConfigReaderFound
}

//ReadMocksDefinition reads all definitions and return an array of valid mocks
func (fd *FileDefinition) ReadMocksDefinition() []Mock {

	if !fd.existsConfigPath(fd.Path) {
		log.Fatalf(ErrNotFoundPath.Error())
	}

	mocks := []Mock{}
	for _, file := range fd.getConfigFiles(fd.Path) {
		if mockDef, err := fd.getMockFromFile(file); err == nil {
			mockDef.Name = filepath.Base(file)
			mocks = append(mocks, mockDef)
		}
	}

	sort.Sort(PrioritySort(mocks))

	return mocks
}
