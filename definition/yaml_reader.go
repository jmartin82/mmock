package definition

import (
	"io/ioutil"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/jmartin82/mmock/logging"
)

//YAMLReader struct created to read yaml config files
type YAMLReader struct {
}

//CanRead return true if is a yaml file
func (jp YAMLReader) CanRead(filename string) bool {
	return filepath.Ext(filename) == ".yaml"
}

//Read Unmarshal a yaml file to Mock struct
func (jp YAMLReader) Read(filename string) (Mock, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return Mock{}, err
	}
	logging.Printf("Loading YAML config: %s\n", filename)
	m := Mock{}
	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		logging.Printf("Invalid mock definition in: %s\n", filename)
		logging.Printf(err.Error())
		return Mock{}, err
	}
	return m, nil
}
