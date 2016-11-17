package definition

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
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
	log.Printf("Loading JSON config: %s\n", filename)
	m := Mock{}
	err = json.Unmarshal(buf, &m)
	if err != nil {
		log.Printf("Invalid mock definition in: %s\n", filename)
		return Mock{}, err
	}
	return m, nil
}
