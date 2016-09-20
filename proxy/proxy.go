package proxy

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jmartin82/mmock/definition"
)

// Proxy calls to real service
type Proxy struct {
	URL string
}

func (pr *Proxy) makeRequest(request definition.Request, extraHeaders definition.Values) definition.Response {

	log.Println("URL:>", pr.URL)

	req, err := http.NewRequest(request.Method, pr.URL, bytes.NewBufferString(request.Body))
	for h, values := range request.Headers {
		for _, value := range values {
			req.Header.Add(h, value)
		}
	}

	for h, values := range extraHeaders {
		for _, value := range values {
			req.Header.Add(h, value)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	r := definition.Response{}
	r.StatusCode = resp.StatusCode
	body, _ := ioutil.ReadAll(resp.Body)
	r.Body = string(body)
	return r
}
