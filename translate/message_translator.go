package translate

import (
	"net/http"

	"github.com/jmartin82/mmock/definition"
)

//MockRequestBuilder defines the translator from http.Request to definition.Request
type MockRequestBuilder interface {
	BuildRequestDefinitionFromHTTP(req *http.Request) definition.Request
}

//MockResponseWriter defines the translator from definition.Response to http.ResponseWriter
type MockResponseWriter interface {
	WriteHTTPResponseFromDefinition(fr *definition.Response, w http.ResponseWriter)
}

//MessageTranslator defines the translator contract between http and mock and viceversa.
//this translation decople the mock matcher from the specific http implementation.
type MessageTranslator interface {
	MockRequestBuilder
	MockResponseWriter
}
