package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jmartin82/mmock/v3/pkg/mock"
)

// HandleCallback makes a callback required after a request
func HandleCallback(cb mock.Callback) (*http.Response, error) {
	if d := cb.Delay.Duration; d > 0 {
		log.Printf("Delaying callback by: %s\n", d)
		time.Sleep(d)
	}

	url := cb.Url
	log.Printf("Making callback to %s\n", url)
	req, err := http.NewRequest(cb.Method, url, bytes.NewBufferString(cb.Body))
	if err != nil {
		return nil, fmt.Errorf("Error creating HTTP request: %w", err)
	}
	// add headers
	for h, vs := range cb.Headers {
		for _, v := range vs {
			req.Header.Set(h, v)
		}
	}

	// Default timeout 10s
	timeout := time.Second * 10
	if cb.Timeout.Duration > 0 {
		timeout = cb.Timeout.Duration
	}
	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making HTTP request: %w", err)
	}
	statusCode := resp.StatusCode
	if statusCode >= 400 {
		body, errBody := ioutil.ReadAll(resp.Body)
		if errBody != nil {
			// Can't read a body, just log the statusCode
			err = fmt.Errorf("Unexpected HTTP response. Status code %d", statusCode)
			return resp, err
		}

		// Include the response body
		err = fmt.Errorf("Unexpected HTTP response. Status code %d, Body %s", statusCode, body)
		return resp, err
	}
	return resp, nil
}
