package definition

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/radovskyb/watcher"
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
	watcher       *watcher.Watcher
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
	fd.watcher = watcher.New()

	// SetMaxEvents to 1 to allow at most 1 Event to be received
	fd.watcher.SetMaxEvents(1)
	go func() {
		for {
			select {
			case event := <-fd.watcher.Event:
				log.Println("Changes detected in mock definitions ", event.String())
				fd.Updates <- fd.ReadMocksDefinition()
			case err := <-fd.watcher.Error:
				log.Println("File monitor error", err)
			}
		}
	}()

	// Watch dir recursively for changes.
	if err := fd.watcher.AddRecursive(fd.Path); err != nil {
		log.Println("Impossible bind the config folder to the files monitor: ", err)
		return
	}

	go func() {
		if err := fd.watcher.Start(time.Millisecond * 100); err != nil {
			log.Println("Impossible to start the config files monitor: ", err)
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
			log.Printf("Loading config file: %s\n", filename)
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
