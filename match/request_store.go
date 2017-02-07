package match

import (
	"github.com/jmartin82/mmock/definition"
)

type RequestStore interface {
	Save(definition.Request)
	Reset()
	GetRequests() []definition.Request
}
