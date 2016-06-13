package translate

import (
	"mmock/definition"
	"net/http"
)

type MockRequestBuilder interface {
	BuildRequestDefinitionFromHTTP(req *http.Request) definition.Request
}
type MockResponseWriter interface {
	WriteHTTPResponseFromDefinition(fr *definition.Response, w http.ResponseWriter)
}

type MessageTranslator interface {
	MockRequestBuilder
	MockResponseWriter
}
