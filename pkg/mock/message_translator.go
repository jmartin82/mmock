package mock

import (
	"net/http"
)

// MockRequestBuilder defines the translator from http.Request to config.Request
type MockRequestBuilder interface {
	BuildRequestDefinitionFromHTTP(req *http.Request) Request
}

// MockResponseWriter defines the translator from config.Response to http.ResponseWriter
type MockResponseWriter interface {
	WriteHTTPResponseFromDefinition(fr *Response, w http.ResponseWriter, req *http.Request)
}

// MessageTranslator defines the translator contract between http and mock and viceversa.
// this translation decople the mock checker from the specific http implementation.
type MessageTranslator interface {
	MockRequestBuilder
	MockResponseWriter
}
