package definition

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
)

//ErrNotValidParserFound we don't have any config reader valid for this file
var ErrNotValidParserFound = errors.New("Not valid config reader found")

//ErrInvalidMockDefinition the file contains an invalid mock definition
var ErrInvalidMockDefinition = errors.New("Invalid mock definition")

//Reader interface contains the funtions to obtain the mock defintions.
type Reader interface {
	Read(string) (Mock, error)
}

//Write interface contains the function persiste mock definitions.
type Writer interface {
	Write(string, Mock) error
}

//Parser interface allows recognize if there is available some config reader for an a specific file.
type Parser interface {
	CanParse(filename string) bool
	Parse(content []byte) (Mock, error)
}

//NewConfigMapper file definition constructor
func NewConfigMapper() *ConfigMapper {
	return &ConfigMapper{
		parsers: []Parser{},
	}
}

//ConfigMapper this struct contains the path of definition and some config readers
type ConfigMapper struct {
	parsers []Parser
}

//AddConfigPaser allows append new readers to able load different config files
func (fd *ConfigMapper) AddConfigParser(reader Parser) {
	fd.parsers = append(fd.parsers, reader)
}

func (fd *ConfigMapper) Write(filename string, mock Mock) error {

	content, err := json.MarshalIndent(mock, "", "  ")
	if err != nil {
		return err
	}

	ioutil.WriteFile(filename, content, 0644)

	return nil
}

func (fd *ConfigMapper) Read(filename string) (Mock, error) {
	for _, parser := range fd.parsers {
		if parser.CanParse(filename) {
			buf, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Printf("Invalid mock definition in: %s\n", filename)
				return Mock{}, ErrInvalidMockDefinition
			}
			log.Printf("Loading config file: %s\n", filename)
			mock, erd := parser.Parse(buf)
			if erd != nil {
				log.Printf("Invalid mock format in: %s Err: %s", filename, erd)
			}
			return mock, erd
		}

	}
	return Mock{}, ErrNotValidParserFound
}
