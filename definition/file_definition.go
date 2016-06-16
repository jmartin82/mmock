package definition

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/fsnotify/fsnotify"
)

var ErrNotFoundPath = errors.New("Configuration path not found")

type FileDefinition struct {
	Path    string
	Updates chan []Mock
}

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

func (fd FileDefinition) existsConfigPath(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (fd FileDefinition) getConfigFiles(path string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	cf := make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) == ".json" {
			cf = append(cf, file.Name())
		}

	}
	return cf
}

func (fd FileDefinition) readMock(filename string) (Mock, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return Mock{}, err
	}
	m := Mock{}
	err = json.Unmarshal(buf, &m)
	if err != nil {
		log.Printf("Invalid mock definition in: %s\n", filename)
		return Mock{}, err
	}
	m.Name = filepath.Base(filename)
	return m, nil
}

func (fd FileDefinition) WatchDir() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("Hot mock file changing not available.")
		return
	}
	err = watcher.Add(fd.Path)
	if err != nil {
		log.Printf("Hot mock file changing not available in folder: %s\n", fd.Path)
		return
	}
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if filepath.Ext(event.Name) == ".json" && (event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Remove == fsnotify.Remove) {

					log.Println("Changes detected in mock definitions")
					fd.Updates <- fd.ReadMocksDefinition()

				}
			case err := <-watcher.Errors:
				log.Fatal(err)
			}
		}
	}()
}

func (fd FileDefinition) ReadMocksDefinition() []Mock {

	if !fd.existsConfigPath(fd.Path) {
		log.Fatalf(ErrNotFoundPath.Error())
	}

	mocks := make([]Mock, 0)

	for _, name := range fd.getConfigFiles(fd.Path) {
		filename := path.Join(fd.Path, name)
		log.Println("Loading mock definition: ", filename)
		if mockDef, err := fd.readMock(filename); err == nil {
			mocks = append(mocks, mockDef)
		}

	}

	sort.Sort(PrioritySort(mocks))

	return mocks
}
