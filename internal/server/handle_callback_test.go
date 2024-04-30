package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmartin82/mmock/v3/pkg/mock"
)

func makeTestServer(allowedMethod string, successResponse string, t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == allowedMethod {
			t.Logf("Returning 200 %s", successResponse)
			w.WriteHeader(200)
			fmt.Fprint(w, successResponse)
		} else {
			t.Logf("Returning 405")
			w.WriteHeader(405)
			fmt.Fprint(w, "Invalid method")
		}
	}))
}

func TestBadStatusCodeErrors(t *testing.T) {
	ts := makeTestServer("GET", "response", t)
	defer ts.Close()

	cb := mock.Callback{
		Url:    ts.URL,
		Method: "POST",
	}

	_, err := HandleCallback(cb)
	if err == nil {
		t.Fatalf("Expected error from HandleCallback, got nil")
	}
}

func TestPost(t *testing.T) {
	postResponse := "Post Response"
	ts := makeTestServer("POST", postResponse, t)
	defer ts.Close()

	cb := mock.Callback{
		Url:        ts.URL,
		Method:     "POST",
		HTTPEntity: mock.HTTPEntity{Body: "Some post body"},
	}

	resp, err := HandleCallback(cb)
	defer resp.Body.Close()
	if err != nil {
		t.Fatalf("Unexpected error from HandleCallback %s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if string(body) != postResponse {
		t.Fatalf("Unexpected response body: %s", body)
	}
}

func TestGet(t *testing.T) {
	getResponse := "Get Response"
	ts := makeTestServer("GET", getResponse, t)
	defer ts.Close()

	cb := mock.Callback{
		Url:    ts.URL,
		Method: "GET",
	}

	resp, err := HandleCallback(cb)
	defer resp.Body.Close()
	if err != nil {
		t.Fatalf("Unexpected error from HandleCallback %s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if string(body) != getResponse {
		t.Fatalf("Unexpected response body: %s", body)
	}
}
