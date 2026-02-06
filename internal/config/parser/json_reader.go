package parser

import (
	"encoding/json"
	"github.com/jmartin82/mmock/v3/pkg/mock"
        "github.com/jmartin82/mmock/v3/internal/config/logger"
	"path/filepath"
	"strings"
)

var log = logger.Log

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
          log.Errorf("JSONReader Parse error: %v", err)
          log.Errorf("JSONReader attempted to parse: %v", string(buf))
          return mock.Definition{}, err
	}

	return m, nil
}
