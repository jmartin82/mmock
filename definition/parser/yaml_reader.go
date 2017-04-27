package parser

import (
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/jmartin82/mmock/definition"
)

//YAMLReader struct created to read yaml config files
type YAMLReader struct {
}

//CanParse return true if is a yaml file
func (jp YAMLReader) CanParse(filename string) bool {
	upperExt := strings.ToUpper(filepath.Ext(filename))
	return ".YAML" == upperExt || ".YML" == upperExt
}

//Read Unmarshal a yaml file to Mock struct
func (jp YAMLReader) Parse(buf []byte) (definition.Mock, error) {
	m := definition.Mock{}
	err := yaml.Unmarshal(buf, &m)
	if err != nil {
		return definition.Mock{}, err
	}
	return m, nil
}
