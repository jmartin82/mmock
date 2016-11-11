package parse

import (
	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/persist"
)

//ResponseParser contains the functions to replace mock response tag
//For instance, it replaces fake.* for some random data or request.* for some provided data in the request.
type ResponseParser interface {
	//Parse subtitutes the current mock response and replace the tags stored inside.
	Parse(*definition.Request, *definition.Response, persist.BodyPersister)
}
