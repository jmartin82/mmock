package config

import "github.com/jmartin82/mmock/pkg/mock"

type Mapping interface {
	Set(URI string, mock mock.Definition) error
	Delete(URI string) error
	Get(URI string) (mock.Definition, bool)
	List() []mock.Definition
}
