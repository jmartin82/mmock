package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// ResponseWriter wraps http.ResponseWriter with gzip compression
type ResponseWriter struct {
	http.ResponseWriter
	Writer     io.Writer
	gzipWriter *gzip.Writer
}

// NewResponseWriter creates gzip wrapper
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	gz := gzip.NewWriter(w)
	return &ResponseWriter{
		ResponseWriter: w,
		Writer:         gz,
		gzipWriter:     gz,
	}
}

// Write compresses data
func (w *ResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// Close flushes gzip stream
func (w *ResponseWriter) Close() error {
	return w.gzipWriter.Close()
}

// ShouldCompress checks if content-type should be compressed
func ShouldCompress(contentType string) bool {
	compressibleTypes := []string{
		"text/",
		"application/json",
		"application/xml",
		"application/javascript",
		"application/x-javascript",
		"application/xhtml+xml",
	}

	contentType = strings.ToLower(contentType)
	for _, t := range compressibleTypes {
		if strings.HasPrefix(contentType, t) {
			return true
		}
	}
	return false
}

// AcceptsGzip checks if request accepts gzip
func AcceptsGzip(req *http.Request) bool {
	encoding := req.Header.Get("Accept-Encoding")
	return strings.Contains(strings.ToLower(encoding), "gzip")
}
