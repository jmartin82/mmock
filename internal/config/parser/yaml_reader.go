package parser

import (
	"github.com/jmartin82/mmock/v3/pkg/mock"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
)

//YAMLReader struct created to read yaml config files
type YAMLReader struct {
}

//CanParse return true if is a yaml file
func (jp YAMLReader) CanParse(filename string) bool {
	upperExt := strings.ToUpper(filepath.Ext(filename))
	return ".YAML" == upperExt || ".YML" == upperExt
}

//Read Unmarshal a yaml file to Definition struct
func (jp YAMLReader) Parse(buf []byte) (mock.Definition, error) {
	m := mock.Definition{}
	err := yaml.Unmarshal(buf, &m)
	if err != nil {
		return mock.Definition{}, err
	}
	return m, nil
}
