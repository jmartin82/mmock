package parser

import (
	"encoding/json"
	"github.com/jmartin82/mmock/v3/pkg/mock"
	"path/filepath"
	"strings"
)

// JSONReader struct created to read json config files
type JSONReader struct {
}

// CanParse return true if is a json file
func (jp JSONReader) CanParse(filename string) bool {
	return ".JSON" == strings.ToUpper(filepath.Ext(filename))
}

// Read Unmarshal a json file to Definition struct
func (jp JSONReader) Parse(buf []byte) (mock.Definition, error) {
	m := mock.Definition{}
	err := json.Unmarshal(buf, &m)
	if err != nil {
		return mock.Definition{}, err
	}
	return m, nil
}
