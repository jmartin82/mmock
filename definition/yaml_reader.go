package definition

import (
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
)

//YAMLReader struct created to read yaml config files
type YAMLReader struct {
}

//CanRead return true if is a yaml file
func (jp YAMLReader) CanRead(filename string) bool {
	upperExt := strings.ToUpper(filepath.Ext(filename))
	return ".YAML" == upperExt || ".YML" == upperExt
}

//Read Unmarshal a yaml file to Mock struct
func (jp YAMLReader) Read(buf []byte) (Mock, error) {
	m := Mock{}
	err := yaml.Unmarshal(buf, &m)
	if err != nil {
		return Mock{}, err
	}
	return m, nil
}
