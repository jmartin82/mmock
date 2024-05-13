package vars

import (
	"strings"
	"testing"

	"github.com/jmartin82/mmock/v3/pkg/mock"
)

var testVars = []struct {
	key          string
	value        string
	expectToFind bool
}{
	{"response.body.test1", "one", true},
	{"response.body.test2", "two", true},
	{"response.body.not_found", "nothing", false},
	{"response.header.Content-Type", "application/json", true},
	{"response.header.not_found", "nothing", false},
	{"response.cookie.test_cookie", "test_cookie_value", true},
	{"response.cookie.not_found", "nothing", false},
}

func getTestResponse() *mock.Response {
	var status = 200

	var headers = make(mock.Values)
	headers["Content-Type"] = []string{"application/json"}

	var cookies = make(mock.Cookies)
	cookies["test_cookie"] = "test_cookie_value"

	var httpHeaders = mock.HttpHeaders{Headers: headers, Cookies: cookies}

	var body = "{\"test1\":\"one\",\"test2\":\"two\"}"
	var httpEntity = mock.HTTPEntity{HttpHeaders: httpHeaders, Body: body}

	var response = mock.Response{StatusCode: status, HTTPEntity: httpEntity}
	return &response
}

func getLoadedResponseFiller() ResponseFiller {
	return ResponseFiller{Response: getTestResponse()}
}

func TestResponseFiller(t *testing.T) {
	var filler = getLoadedResponseFiller()

	for _, tt := range testVars {
		holders := []string{
			tt.key,
		}

		vars := filler.Fill(holders)

		if len(vars) == 0 {
			if tt.expectToFind {
				t.Errorf("Unable to retrieve vars")
			}
			continue
		}

		v, ok := vars[tt.key]

		if !ok {
			t.Errorf("Unable to retrieve value inside vars")
			continue
		}

		if strings.EqualFold(v[0], tt.value) != tt.expectToFind {
			t.Errorf("Couldn't get the expected value. Expected: %s, Value found: %s", tt.value, v[0])
		}
	}
}
