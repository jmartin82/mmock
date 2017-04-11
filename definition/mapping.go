package definition

import (
	"errors"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

var ErrFilePathIsNotUnderConfigPath = errors.New("File path is not under config path")

type Mapping interface {
	Set(URI string, mock Mock) error
	Delete(URI string) error
	List() []Mock
	Load(URI string) error
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

type ConfigMapping struct {
	mapper  *ConfigMapper
	mapping map[string]Mock
	path    string
	sync.Mutex
}

func NewConfigMapping(path string, mapper *ConfigMapper) *ConfigMapping {
	return &ConfigMapping{path: path, mapper: mapper, mapping: make(map[string]Mock)}
}

func (fm *ConfigMapping) Set(URI string, mock Mock) error {
	defer fm.Unlock()
	fm.Lock()
	fileName, err := fm.resolveFile(URI)
	if err != nil {
		return err
	}

	if err := fm.mapper.Write(fileName, mock); err != nil {
		return err
	}

	fm.mapping[URI] = mock
	return nil
}
func (fm *ConfigMapping) Delete(URI string) error {
	defer fm.Unlock()
	fm.Lock()
	fileName, err := fm.resolveFile(URI)
	if err != nil {
		return err
	}

	if err := os.Remove(fileName); err != nil {
		return err
	}

	delete(fm.mapping, URI)

	return nil
}
func (fm *ConfigMapping) Load(URI string) error {
	defer fm.Unlock()
	fm.Lock()
	fileName, errf := fm.resolveFile(URI)
	if errf != nil {
		return errf
	}

	mock, err := fm.mapper.Read(fileName)
	mock.Name = URI
	if err != nil {
		return err

	}

	fm.mapping[URI] = mock

	return nil
}

func (fm *ConfigMapping) List() []Mock {
	defer fm.Unlock()
	fm.Lock()
	mocks := make([]Mock, len(fm.mapping))
	for _, mock := range fm.mapping {
		mocks = append(mocks, mock)
	}
	sort.Sort(PrioritySort(mocks))

	return mocks
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
