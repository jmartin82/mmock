package route

import (
	"github.com/jmartin82/mmock/definition"
)

//Router contains the functions to check the http request and return the matching mock.
type Router interface {
	//Route checks the request with all available mock definitions and return the matching mock for it.
	Route(req *definition.Request) (*definition.Mock, map[string]string)
	//SetMockDefinitions allows replace the current mock definitions for new ones.
	SetMockDefinitions(mocks []definition.Mock)
}
