package definition

import (
	"encoding/json"
	"path/filepath"
	"strings"
)

//JSONReader struct created to read json config files
type JSONReader struct {
}

//CanRead return true if is a json file
func (jp JSONReader) CanRead(filename string) bool {
	return ".JSON" == strings.ToUpper(filepath.Ext(filename))
}

//Read Unmarshal a json file to Mock struct
func (jp JSONReader) Read(buf []byte) (Mock, error) {
	m := Mock{}
	err := json.Unmarshal(buf, &m)
	if err != nil {
		return Mock{}, err
	}
	return m, nil
}
