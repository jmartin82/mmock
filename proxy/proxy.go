package proxy

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jmartin82/mmock/definition"
)

// Proxy calls to real service
type Proxy struct {
	URL string
}

// MakeRequest creates a real request to the desired service using data from the original request
func (pr *Proxy) MakeRequest(request definition.Request) *definition.Response {

	r := &definition.Response{}
	log.Println("Proxy to URL:>", pr.URL)
	req, err := http.NewRequest(request.Method, pr.URL, bytes.NewBufferString(request.Body))
	for h, values := range request.Headers {
		for _, value := range values {
			req.Header.Add(h, value)
		}
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Impossible create a proxy request: ", err)
		r.StatusCode = 500
		return r
	}
	defer resp.Body.Close()

	r.StatusCode = resp.StatusCode

	r.Headers = make(definition.Values)
	for h, values := range resp.Header {
		r.Headers[h] = values
	}
	r.Cookies = make(definition.Cookies)
	for _, cookie := range resp.Cookies() {
		r.Cookies[cookie.Name] = cookie.Value
	}

	body, _ := ioutil.ReadAll(resp.Body)
	r.Body = string(body)
	return r
}
