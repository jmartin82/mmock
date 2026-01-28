package gzip

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShouldCompress(t *testing.T) {
	tests := []struct {
		contentType string
		expected    bool
	}{
		{"application/json", true},
		{"Application/JSON", true}, // Case insensitive
		{"text/html", true},
		{"text/plain", true},
		{"text/css", true},
		{"application/xml", true},
		{"application/javascript", true},
		{"application/x-javascript", true},
		{"application/xhtml+xml", true},
		{"image/png", false},
		{"image/jpeg", false},
		{"video/mp4", false},
		{"application/pdf", false},
		{"application/zip", false},
		{"application/gzip", false},
		{"application/octet-stream", false},
		{"", false},
	}

	for _, test := range tests {
		result := ShouldCompress(test.contentType)
		if result != test.expected {
			t.Errorf("ShouldCompress(%q) = %v, want %v",
				test.contentType, result, test.expected)
		}
	}
}

func TestAcceptsGzip(t *testing.T) {
	tests := []struct {
		encoding string
		expected bool
	}{
		{"gzip", true},
		{"gzip, deflate", true},
		{"deflate, gzip", true},
		{"deflate, gzip, br", true},
		{"GZIP", true}, // Case insensitive
		{"deflate", false},
		{"br", false},
		{"", false},
	}

	for _, test := range tests {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Accept-Encoding", test.encoding)

		result := AcceptsGzip(req)
		if result != test.expected {
			t.Errorf("AcceptsGzip(%q) = %v, want %v",
				test.encoding, result, test.expected)
		}
	}
}

func TestGzipResponseWriter(t *testing.T) {
	// Test that data is actually compressed
	recorder := httptest.NewRecorder()
	writer := NewResponseWriter(recorder)

	testData := strings.Repeat("test data ", 100) // Repetitive data compresses well
	writer.Write([]byte(testData))
	writer.Close()

	compressed := recorder.Body.Bytes()

	// Verify data was compressed (should be smaller)
	if len(compressed) >= len(testData) {
		t.Error("Data was not compressed")
	}

	// Verify can decompress
	reader, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		t.Fatalf("Failed to create gzip reader: %v", err)
	}
	defer reader.Close()

	decompressed, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatalf("Failed to decompress: %v", err)
	}

	if string(decompressed) != testData {
		t.Error("Decompressed data doesn't match original")
	}
}

func TestGzipResponseWriterHeader(t *testing.T) {
	// Test that ResponseWriter properly wraps Header() method
	recorder := httptest.NewRecorder()
	writer := NewResponseWriter(recorder)

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("X-Custom-Header", "test-value")

	if recorder.Header().Get("Content-Type") != "application/json" {
		t.Error("Header() not properly wrapped")
	}

	if recorder.Header().Get("X-Custom-Header") != "test-value" {
		t.Error("Custom headers not properly wrapped")
	}
}

func TestGzipResponseWriterWriteHeader(t *testing.T) {
	// Test that ResponseWriter properly wraps WriteHeader() method
	recorder := httptest.NewRecorder()
	writer := NewResponseWriter(recorder)

	writer.WriteHeader(http.StatusCreated)

	if recorder.Code != http.StatusCreated {
		t.Errorf("WriteHeader() not properly wrapped: got %d, want %d",
			recorder.Code, http.StatusCreated)
	}
}
