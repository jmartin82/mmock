package parse

import "github.com/jmartin82/mmock/definition"

//ResponseParser contains the functions to replace mock response tag
//For instance, it replaces fake.* for some random data or request.* for some provided data in the request.
type ResponseParser interface {
	//Parse subtitutes the current mock response and replace the tags stored inside.
	Parse(*definition.Request, *definition.Response)
	//ReplaceVars relplaces variables from the request in the input
	ReplaceVars(req *definition.Request, input string) string
	//ParseBody parses body respecting bodyAppend and replacing variables from request
	ParseBody(req *definition.Request, body string, bodyAppend string) string
}
