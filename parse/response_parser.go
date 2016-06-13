package parse

import (
	"github.com/jmartin82/mmock/definition"
)

type ResponseParser interface {
	Parse(*definition.Request, *definition.Response)
}
