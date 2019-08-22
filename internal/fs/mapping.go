package fs

import (
	"errors"
	"github.com/jmartin82/mmock/pkg/mock"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

var ErrFilePathIsNotUnderConfigPath = errors.New("File path is not under config path")
var ErrMockDoesntExist = errors.New("Definition doesn't exist")

//PrioritySort mock array sorted by priority
type PrioritySort []mock.Definition

func (s PrioritySort) Len() int {
	return len(s)
}
func (s PrioritySort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s PrioritySort) Less(i, j int) bool {
	return s[i].Control.Priority > s[j].Control.Priority
}

type ConfigMapping struct {
	mapper   *FSMapper
	mapping  map[string]mock.Definition
	path     string
	fsListen bool
	fsUpdate chan struct{}
	sync.Mutex
}

func NewConfigMapping(path string, mapper *FSMapper, fsUpdate chan struct{}) *ConfigMapping {
	cm := &ConfigMapping{path: path, mapper: mapper, mapping: make(map[string]mock.Definition), fsUpdate: fsUpdate}
	cm.populate()
	cm.fsBind()
	go cm.listenFsChanges()
	return cm
}

func (fm *ConfigMapping) listenFsChanges() {
	for {
		<-fm.fsUpdate
		if fm.fsIsBind() {
			fm.populate()
		}

	}
}

func (fm *ConfigMapping) Get(URI string) (mock.Definition, bool) {
	URI = fm.sanitizeURI(URI)
	mock, ok := fm.mapping[URI]
	return mock, ok
}

func (fm *ConfigMapping) Set(URI string, mock mock.Definition) error {
	defer fm.Unlock()
	fm.Lock()
	fm.fsUnBind()
	URI = fm.sanitizeURI(URI)
	fileName, err := fm.resolveFile(URI)
	if err != nil {
		return err
	}

	if err := fm.mapper.Write(fileName, mock); err != nil {
		return err
	}

	fm.mapping[URI] = mock
	fm.fsBind()
	return nil
}
func (fm *ConfigMapping) Delete(URI string) error {

	defer fm.Unlock()
	fm.Lock()
	fm.fsUnBind()
	URI = fm.sanitizeURI(URI)
	fileName, err := fm.resolveFile(URI)
	if err != nil {
		return err
	}

	if err := os.Remove(fileName); err != nil {
		return err
	}

	delete(fm.mapping, URI)
	fm.fsBind()
	return nil
}

func (fm *ConfigMapping) List() []mock.Definition {
	defer fm.Unlock()
	fm.Lock()
	mocks := make([]mock.Definition, 0, len(fm.mapping))
	for _, mock := range fm.mapping {

		mocks = append(mocks, mock)
	}

	sort.Sort(PrioritySort(mocks))

	return mocks
}

func (fm *ConfigMapping) populate() {
	fm.mapping = make(map[string]mock.Definition)
	filepath.Walk(fm.path, func(filePath string, fileInfo os.FileInfo, err error) error {
		if !fileInfo.IsDir() {
			URI := strings.TrimPrefix(filePath, fm.path)
			if err := fm.load(URI); err != nil {
				log.Printf("Error %v. Loading config: %v\n", err, URI)
			}
		}
		return nil
	})
}

func (fm *ConfigMapping) load(URI string) error {
	defer fm.Unlock()
	fm.Lock()
	URI = fm.sanitizeURI(URI)
	fileName, errf := fm.resolveFile(URI)
	if errf != nil {
		return errf
	}

	mock, err := fm.mapper.Read(fileName)
	mock.URI = URI
	if err != nil {
		return err

	}

	fm.mapping[URI] = mock

	return nil
}

func (fm *ConfigMapping) resolveFile(URI string) (string, error) {
	filename, err := filepath.Abs(path.Join(fm.path, URI))
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(filename, fm.path) {
		log.Printf("File path not under the config path\n")
		return "", ErrFilePathIsNotUnderConfigPath
	}
	return filename, nil
}

func (fm *ConfigMapping) fsUnBind() {
	fm.fsListen = false
}

func (fm *ConfigMapping) fsBind() {
	fm.fsListen = true
}

func (fm *ConfigMapping) fsIsBind() bool {
	return fm.fsListen
}

func (fm *ConfigMapping) sanitizeURI(URI string) string {
	return strings.Trim(strings.TrimPrefix(URI, string(os.PathSeparator)), " ")
}
