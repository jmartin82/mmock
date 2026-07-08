package server

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmartin82/mmock/v3/pkg/match"
	"github.com/jmartin82/mmock/v3/pkg/mock"
)

// mockTranslator is a simple mock implementation for testing
type mockTranslator struct{}

func (mt mockTranslator) BuildRequestDefinitionFromHTTP(req *http.Request) mock.Request {
	return mock.Request{
		Method: req.Method,
		Path:   req.URL.Path,
	}
}

func (mt mockTranslator) WriteHTTPResponseFromDefinition(fr *mock.Response, w http.ResponseWriter) {
	// Write headers
	for header, values := range fr.Headers {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}

	// Write status code
	w.WriteHeader(fr.StatusCode)

	// Write body
	w.Write([]byte(fr.Body))
}

// mockResolver returns predefined mocks for testing
type mockResolver struct {
	response *mock.Response
}

func (mr *mockResolver) Resolve(req *mock.Request) (*mock.Definition, *match.Result) {
	def := &mock.Definition{
		Response: *mr.response,
		Control:  mock.Control{},
	}
	result := &match.Result{
		Found: true,
		URI:   "test-mock",
	}
	return def, result
}

// mockSpier does nothing for testing
type mockSpier struct{}

func (ms *mockSpier) Save(t match.Transaction)                      {}
func (ms *mockSpier) Reset()                                        {}
func (ms *mockSpier) ResetMatch(r mock.Request)                     {}
func (ms *mockSpier) GetAll() []match.Transaction                   { return nil }
func (ms *mockSpier) Get(limit int, offset int) []match.Transaction { return nil }
func (ms *mockSpier) Find(r mock.Request) []match.Transaction       { return nil }
func (ms *mockSpier) GetMatched() []match.Transaction               { return nil }
func (ms *mockSpier) GetUnMatched() []match.Transaction             { return nil }

// mockEvaluator does nothing for testing
type mockEvaluator struct{}

func (me *mockEvaluator) Eval(req *mock.Request, m *mock.Definition) {}

// mockScenario does nothing for testing
type mockScenario struct{}

func (ms *mockScenario) SetState(name, status string) {}
func (ms *mockScenario) GetState(name string) string  { return "not_started" }
func (ms *mockScenario) Reset(name string) bool       { return false }
func (ms *mockScenario) ResetAll()                    {}
func (ms *mockScenario) SetPaused(newstate bool)      {}
func (ms *mockScenario) GetPaused() bool              { return false }
func (ms *mockScenario) List() string                 { return "{}" }

func TestDispatcherGzipJSON(t *testing.T) {
	// Setup mock response with JSON content type
	response := &mock.Response{
		StatusCode: 200,
		HTTPEntity: mock.HTTPEntity{
			HttpHeaders: mock.HttpHeaders{
				Headers: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			Body: `{"message": "test data"}`,
		},
	}

	dispatcher := &Dispatcher{
		Translator: mockTranslator{},
		Resolver:   &mockResolver{response: response},
		Evaluator:  &mockEvaluator{},
		Scenario:   &mockScenario{},
		Spier:      &mockSpier{},
		Mlog:       make(chan match.Transaction, 1),
	}

	// Create request with Accept-Encoding: gzip
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	recorder := httptest.NewRecorder()

	// Serve HTTP
	dispatcher.ServeHTTP(recorder, req)

	// Verify Content-Encoding header
	if recorder.Header().Get("Content-Encoding") != "gzip" {
		t.Error("Content-Encoding header not set to gzip")
	}

	// Verify response is compressed
	reader, err := gzip.NewReader(recorder.Body)
	if err != nil {
		t.Fatalf("Response not gzip compressed: %v", err)
	}
	defer reader.Close()

	body, _ := ioutil.ReadAll(reader)
	if string(body) != response.Body {
		t.Errorf("Got %s, want %s", body, response.Body)
	}
}

func TestDispatcherNoGzipForImages(t *testing.T) {
	// Setup mock response with image content type
	response := &mock.Response{
		StatusCode: 200,
		HTTPEntity: mock.HTTPEntity{
			HttpHeaders: mock.HttpHeaders{
				Headers: map[string][]string{
					"Content-Type": {"image/png"},
				},
			},
			Body: "binary image data",
		},
	}

	dispatcher := &Dispatcher{
		Translator: mockTranslator{},
		Resolver:   &mockResolver{response: response},
		Evaluator:  &mockEvaluator{},
		Scenario:   &mockScenario{},
		Spier:      &mockSpier{},
		Mlog:       make(chan match.Transaction, 1),
	}

	// Create request with Accept-Encoding: gzip
	req := httptest.NewRequest("GET", "/image.png", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	recorder := httptest.NewRecorder()

	// Serve HTTP
	dispatcher.ServeHTTP(recorder, req)

	// Should NOT have Content-Encoding header
	if recorder.Header().Get("Content-Encoding") == "gzip" {
		t.Error("Image should not be gzip compressed")
	}

	// Body should be uncompressed
	if recorder.Body.String() != response.Body {
		t.Error("Image body was modified")
	}
}

func TestDispatcherNoGzipWithoutAcceptEncoding(t *testing.T) {
	// Setup mock response with JSON content type
	response := &mock.Response{
		StatusCode: 200,
		HTTPEntity: mock.HTTPEntity{
			HttpHeaders: mock.HttpHeaders{
				Headers: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			Body: `{"message": "test data"}`,
		},
	}

	dispatcher := &Dispatcher{
		Translator: mockTranslator{},
		Resolver:   &mockResolver{response: response},
		Evaluator:  &mockEvaluator{},
		Scenario:   &mockScenario{},
		Spier:      &mockSpier{},
		Mlog:       make(chan match.Transaction, 1),
	}

	// Create request WITHOUT Accept-Encoding header
	req := httptest.NewRequest("GET", "/test", nil)

	recorder := httptest.NewRecorder()

	// Serve HTTP
	dispatcher.ServeHTTP(recorder, req)

	// Should NOT have Content-Encoding header
	if recorder.Header().Get("Content-Encoding") == "gzip" {
		t.Error("Response should not be compressed without Accept-Encoding")
	}

	// Body should be uncompressed
	if recorder.Body.String() != response.Body {
		t.Error("Response body was modified")
	}
}

func TestGetContentType(t *testing.T) {
	tests := []struct {
		name     string
		response *mock.Response
		expected string
	}{
		{
			name: "JSON content type",
			response: &mock.Response{
				HTTPEntity: mock.HTTPEntity{
					HttpHeaders: mock.HttpHeaders{
						Headers: map[string][]string{
							"Content-Type": {"application/json"},
						},
					},
				},
			},
			expected: "application/json",
		},
		{
			name: "No content type",
			response: &mock.Response{
				HTTPEntity: mock.HTTPEntity{
					HttpHeaders: mock.HttpHeaders{
						Headers: map[string][]string{},
					},
				},
			},
			expected: "text/plain",
		},
		{
			name: "Empty content type array",
			response: &mock.Response{
				HTTPEntity: mock.HTTPEntity{
					HttpHeaders: mock.HttpHeaders{
						Headers: map[string][]string{
							"Content-Type": {},
						},
					},
				},
			},
			expected: "text/plain",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := getContentType(test.response)
			if result != test.expected {
				t.Errorf("getContentType() = %s, want %s", result, test.expected)
			}
		})
	}
}

func TestDispatcherGzipTextHTML(t *testing.T) {
	// Setup mock response with HTML content type
	response := &mock.Response{
		StatusCode: 200,
		HTTPEntity: mock.HTTPEntity{
			HttpHeaders: mock.HttpHeaders{
				Headers: map[string][]string{
					"Content-Type": {"text/html"},
				},
			},
			Body: "<html><body>test</body></html>",
		},
	}

	dispatcher := &Dispatcher{
		Translator: mockTranslator{},
		Resolver:   &mockResolver{response: response},
		Evaluator:  &mockEvaluator{},
		Scenario:   &mockScenario{},
		Spier:      &mockSpier{},
		Mlog:       make(chan match.Transaction, 1),
	}

	req := httptest.NewRequest("GET", "/test.html", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	recorder := httptest.NewRecorder()

	dispatcher.ServeHTTP(recorder, req)

	// Verify Content-Encoding header
	if recorder.Header().Get("Content-Encoding") != "gzip" {
		t.Error("HTML content should be gzip compressed")
	}

	// Verify response can be decompressed
	reader, err := gzip.NewReader(recorder.Body)
	if err != nil {
		t.Fatalf("Response not gzip compressed: %v", err)
	}
	defer reader.Close()

	body, _ := ioutil.ReadAll(reader)
	if string(body) != response.Body {
		t.Errorf("Got %s, want %s", body, response.Body)
	}
}

func TestDispatcherGzipDeflateEncoding(t *testing.T) {
	// Test with "deflate, gzip" encoding
	response := &mock.Response{
		StatusCode: 200,
		HTTPEntity: mock.HTTPEntity{
			HttpHeaders: mock.HttpHeaders{
				Headers: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			Body: `{"test": "data"}`,
		},
	}

	dispatcher := &Dispatcher{
		Translator: mockTranslator{},
		Resolver:   &mockResolver{response: response},
		Evaluator:  &mockEvaluator{},
		Scenario:   &mockScenario{},
		Spier:      &mockSpier{},
		Mlog:       make(chan match.Transaction, 1),
	}

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "deflate, gzip")

	recorder := httptest.NewRecorder()

	dispatcher.ServeHTTP(recorder, req)

	// Should still compress with gzip
	if recorder.Header().Get("Content-Encoding") != "gzip" {
		t.Error("Should compress when Accept-Encoding contains gzip")
	}

	// Verify can decompress
	reader, err := gzip.NewReader(bytes.NewReader(recorder.Body.Bytes()))
	if err != nil {
		t.Fatalf("Failed to decompress: %v", err)
	}
	defer reader.Close()

	body, _ := ioutil.ReadAll(reader)
	if string(body) != response.Body {
		t.Errorf("Decompressed body mismatch")
	}
}
