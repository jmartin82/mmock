package config

import (
	"encoding/json"
	"errors"
	"github.com/jmartin82/mmock/v3/pkg/mock"
	"io/ioutil"
	"reflect"
	"strings"
)

// ErrNotValidParserFound we don't have any config reader valid for this file
var ErrNotValidParserFound = errors.New("Not valid config reader found")

// ErrInvalidMockDefinition the file contains an invalid mock config
var ErrInvalidMockDefinition = errors.New("Invalid mock config")

// Reader interface contains the funtions to obtain the mock defintions.
type Reader interface {
	Read(string) (mock.Definition, error)
}

// Writer interface contains the function persist mock definitions.
type Writer interface {
	Write(string, mock.Definition) error
}

// Parser interface allows recognize if there is available some config reader for an a specific file.
type Parser interface {
	CanParse(filename string) bool
	Parse(content []byte) (mock.Definition, error)
}

// NewFileSystemMapper file config constructor
func NewFileSystemMapper() *FSMapper {
	return &FSMapper{
		parsers: []Parser{},
	}
}

// FSMapper this struct contains the path of config and some config readers
type FSMapper struct {
	parsers []Parser
}

// AddParser allows append new readers to able load different config files
func (fd *FSMapper) AddParser(reader Parser) {
	fd.parsers = append(fd.parsers, reader)
}

func (fd *FSMapper) Write(filename string, mock mock.Definition) error {

	content, err := json.MarshalIndent(mock, "", "  ")
	if err != nil {
		return err
	}

	ioutil.WriteFile(filename, content, 0644)

	return nil
}

func (fd *FSMapper) Read(filename string) (mock.Definition, error) {
	for _, parser := range fd.parsers {
		if parser.CanParse(filename) {
			buf, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Errorf("Invalid mock config in: %s\n", filename)
				return mock.Definition{}, ErrInvalidMockDefinition
			}
			log.Infof("Loading config file: %s\n", filename)
			mock, erd := parser.Parse(buf)
			if erd != nil {
				log.Errorf("Invalid mock format in: %s Err: %s", filename, erd)
			}
			if reflect.TypeOf(parser).String() == "parser.YAMLReader" {
				if mock.Request.Body != "" {
					mock.Request.Body = strings.TrimRight(mock.Request.Body, "\n")
				}
			}
			return mock, erd
		}

	}
	return mock.Definition{}, ErrNotValidParserFound
}
