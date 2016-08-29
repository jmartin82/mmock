package definition

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
)

//YAMLReader struct created to read yaml config files
type YAMLReader struct {
}

//CanRead return true if is a yaml file
func (jp YAMLReader) CanRead(filename os.FileInfo) bool {
	return filepath.Ext(filename.Name()) == ".yaml"
}

//Read Unmarshal a yaml file to Mock struct
func (jp YAMLReader) Read(filename string) (Mock, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return Mock{}, err
	}
	log.Printf("Loading YAML config: %s\n", filename)
	m := Mock{}
	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		log.Printf("Invalid mock definition in: %s\n", filename)
		log.Printf(err.Error())
		return Mock{}, err
	}
	return m, nil
}
