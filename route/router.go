package route

import (
	"github.com/jmartin82/mmock/definition"
)

//Router contains the functions to check the http request and return the matching mock.
type Router interface {
	Route(req *definition.Request) (*definition.Mock, definition.MatchErrors)
	SetMockDefinitions(mocks []definition.Mock)
}
