package parser

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/jmartin82/mmock/definition"
)

//JSONReader struct created to read json config files
type JSONReader struct {
}

//CanParse return true if is a json file
func (jp JSONReader) CanParse(filename string) bool {
	return ".JSON" == strings.ToUpper(filepath.Ext(filename))
}

//Read Unmarshal a json file to Mock struct
func (jp JSONReader) Parse(buf []byte) (definition.Mock, error) {
	m := definition.Mock{}
	err := json.Unmarshal(buf, &m)
	if err != nil {
		return definition.Mock{}, err
	}
	return m, nil
}
