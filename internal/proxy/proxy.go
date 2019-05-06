package proxy

import (
	"bytes"
	"github.com/jmartin82/mmock/pkg/mock"
	"io/ioutil"
	"log"
	"net/http"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Proxy calls to real service
type Proxy struct {
	URL    string
	Client HttpClient
}

// MakeRequest creates a real request to the desired service using data from the original request
func (pr *Proxy) MakeRequest(request *mock.Request) *mock.Response {

	r := &mock.Response{}
	log.Println("Proxy to URL:>", pr.URL)
	req, err := http.NewRequest(request.Method, pr.URL, bytes.NewBufferString(request.Body))
	for h, values := range request.Headers {
		for _, value := range values {
			req.Header.Add(h, value)
		}
	}

	q := req.URL.Query()
	for h, values := range request.QueryStringParameters {
		for _, value := range values {
			q.Add(h, value)
		}
	}
	req.URL.RawQuery = q.Encode()

	log.Println("Query string parameters: ", req.URL.RawQuery)
	log.Println("Request body: ", req.Body)

	resp, err := pr.Client.Do(req)
	if err != nil {
		log.Println("Impossible create a proxy request: ", err)
		r.StatusCode = http.StatusInternalServerError
		return r
	}
	defer resp.Body.Close()

	r.StatusCode = resp.StatusCode

	r.Headers = make(mock.Values)
	for h, values := range resp.Header {
		r.Headers[h] = values
	}
	r.Cookies = make(mock.Cookies)
	for _, cookie := range resp.Cookies() {
		r.Cookies[cookie.Name] = cookie.Value
	}

	body, _ := ioutil.ReadAll(resp.Body)
	r.Body = string(body)
	return r
}
