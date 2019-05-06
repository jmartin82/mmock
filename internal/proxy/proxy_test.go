package proxy

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jmartin82/mmock/pkg/mock"
	"io"
	"net/http"
	"testing"
)

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	// just in case you want default correct return value
	return &http.Response{}, nil
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func TestMakeValidRequest(t *testing.T) {

	client := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
			}
			resp.Header = make(http.Header)
			resp.Header.Add("Set-Cookie", "login=1")
			resp.Header.Add("Accept", "text/html")
			resp.Body = nopCloser{bytes.NewBufferString("body data")}
			return resp, nil
		},
	}

	request := &mock.Request{}
	request.Host = "http://mock_host.com"
	request.Method = "GET"
	request.Path = "/home"
	url := "http://example.com"
	proxy := Proxy{URL: url, Client: client}
	response := proxy.MakeRequest(request)

	fmt.Printf("%+v\n", response)

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected code: %d, got: %d", http.StatusOK, response.StatusCode)
	}

	if response.HttpHeaders.Headers["Accept"][0] != "text/html" {
		t.Errorf("Expected accept header == 'text/html'")
	}

	if response.HttpHeaders.Cookies["login"] != "1" {
		t.Errorf("Expected cookie login=1")
	}

	if response.Body != "body data" {
		t.Errorf("Unexpected body")
	}

}

func TestMakeInvalidRequest(t *testing.T) {
	client := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// do whatever you want
			return nil, errors.New("something goes wrong")
		},
	}
	request := &mock.Request{}
	request.Host = "http://mock_host.com"
	request.Method = "GET"
	request.Path = "/home"
	url := "http://example.com"
	proxy := Proxy{URL: url, Client: client}
	response := proxy.MakeRequest(request)
	if response.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected error code: %d, got: %d", http.StatusInternalServerError, response.StatusCode)
	}
}
