package definition

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/jmartin82/mmock/logging"
)

//JSONReader struct created to read json config files
type JSONReader struct {
}

//CanRead return true if is a json file
func (jp JSONReader) CanRead(filename string) bool {
	return filepath.Ext(filename) == ".json"
}

//Read Unmarshal a json file to Mock struct
func (jp JSONReader) Read(filename string) (Mock, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return Mock{}, err
	}
	logging.Printf("Loading JSON config: %s\n", filename)
	m := Mock{}
	err = json.Unmarshal(buf, &m)
	if err != nil {
		logging.Printf("Invalid mock definition in: %s\n", filename)
		return Mock{}, err
	}
	return m, nil
}
