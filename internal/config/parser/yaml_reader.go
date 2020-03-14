package parser

import (
	"bytes"
	"io"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	yamlv2 "gopkg.in/yaml.v2"

	"github.com/jmartin82/mmock/pkg/mock"
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
func (jp YAMLReader) Parse(buf []byte) ([]mock.Definition, error) {
	var mocks []mock.Definition

	dec := yamlv2.NewDecoder(bytes.NewReader(buf))
	for {
		var value interface{}
		if err := dec.Decode(&value); err != nil {
			if err == io.EOF {
				return mocks, nil
			}
			return nil, err
		}

		y, err := yamlv2.Marshal(value)
		if err != nil {
			return nil, err
		}

		m := mock.Definition{}
		if err := yaml.Unmarshal(y, &m); err != nil {
			return nil, err
		}
		mocks = append(mocks, m)
	}

	return mocks, nil
}
