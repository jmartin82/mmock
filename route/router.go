package route

import (
	"github.com/jmartin82/mmock/definition"
)

//Router checks the request with all available mock definitions and return the matching mock for it.
//And also allows replace the current mock definitions for new ones.
type Router interface {
	Route(req *definition.Request) (*definition.Mock, map[string]string)
	SetMockDefinitions(mocks []definition.Mock)
}
