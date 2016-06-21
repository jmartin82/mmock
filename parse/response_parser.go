package parse

import (
	"github.com/jmartin82/mmock/definition"
)

//Response Parser parses the current mock response and replace the tags stored inside.
//For instance, it replaces fake.* for some random data or request.* for some provided data in the request.
type ResponseParser interface {
	Parse(*definition.Request, *definition.Response)
}
