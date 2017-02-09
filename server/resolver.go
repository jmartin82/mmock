package server

import (
	"github.com/jmartin82/mmock/definition"
)

//Resolver contains the functions to check the http request and return the matching mock.
type Resolver interface {
	Resolve(req *definition.Request) (*definition.Mock, definition.MatchErrors)
	SetMockDefinitions(mocks []definition.Mock)
}
