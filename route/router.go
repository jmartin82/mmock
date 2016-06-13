package route

import (
	"github.com/jmartin82/mmock/definition"
)

type Router interface {
	Route(req *definition.Request) (*definition.Mock, map[string]string)
	SetMockDefinitions(mocks []definition.Mock)
}
