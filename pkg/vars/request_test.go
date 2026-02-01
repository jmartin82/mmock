package vars

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/jmartin82/mmock/v3/pkg/mock"
)

func TestGetHeaderParam(t *testing.T) {
	rp := Request{}
	header := make(mock.Values)
	header["Authorization"] = []string{"Bearer abc123"}
	req := mock.Request{}
	req.Headers = header
	rp.Request = &req
	v, f := rp.getHeaderParam("Authorization")
	if !f {
		t.Errorf("Header key not found")
	}

	if !strings.EqualFold(v, "Bearer abc123") {
		t.Errorf("Couldn't get the content. Value: %s", v)
	}
}

func TestGetHeaderParamNotFoundHeaderKey(t *testing.T) {
	rp := Request{}
	header := make(mock.Values)
	header["Authorization"] = []string{"Bearer abc123"}
	req := mock.Request{}
	req.Headers = header
	rp.Request = &req
	_, f := rp.getHeaderParam("Authorization2")
	if f {
		t.Errorf("Header key found")
	}
}

func TestGetHeaderParamWithOutHeaderValue(t *testing.T) {
	rp := Request{}
	header := make(mock.Values)
	header["Authorization"] = []string{}
	req := mock.Request{}
	req.Headers = header
	rp.Request = &req
	v, f := rp.getHeaderParam("Authorization")
	if f {
		t.Errorf("Header key found")
	}

	if strings.EqualFold(v, "Bearer abc1235") {
		t.Errorf("Couldn get the content. Value: %s", v)
	}
}

func TestGetEnvironmentVariable(t *testing.T) {
	key := "env.TEST_VARIABLE"
	value := "test value"

	holders := []string{
		key,
	}

	rp := Request{}

	os.Setenv("TEST_VARIABLE", value)
	vars := rp.Fill(holders)
	os.Unsetenv("TEST_VARIABLE")

	if len(vars) == 0 {
		t.Errorf("Unable to retrieve vars")
	}

	v, ok := vars[key]

	if !ok {
		t.Errorf("Unable to retrieve value inside vars")
	}

	if len(v[0]) == 0 {
		t.Errorf("Unable to retrieve value inside var element")
	}

	if !strings.EqualFold(v[0], value) {
		t.Errorf("Couldn't get the expected value. Expected: %s, Value found: %s", value, v[0])
	}
}

func TestGetEnvironmentVariableNotFound(t *testing.T) {
	key := "env.TEST_VARIABLE_UNEXISTENT"
	value := "test value"

	holders := []string{
		key,
	}

	rp := Request{}

	os.Setenv("TEST_VARIABLE", value)
	vars := rp.Fill(holders)
	os.Unsetenv("TEST_VARIABLE")

	if len(vars) > 0 {
		t.Errorf("Unexpectedly found variable")
	}
}

func TestGetBodyParam(t *testing.T) {
	req := mock.Request{}
	req.Headers = make(mock.Values)
	req.Headers["Content-Type"] = []string{"application/json"}
	req.Body = `
{
  "email": "hilari@sapo.pt",
  "age": 4,
  "uuid":"0bd74115-2307-458f-8288-b726724045ef",
  "discarded": "do not return"
}
`
	res := mock.Response{}
	res.Body = `
{
  "email": "{{request.body.email.regex((\@gmail.com))}}",
  "age": {{request.body.age}},
  "uuid": "{{request.body.uuid.regex(\b([0-9a-zA-Z]{4})\b).concat(-878787)}}",
  "discarded": "{{request.body.discarded.concat(, Please!)}}"
}
`

	expected := `
{
  "email": "",
  "age": 4,
  "uuid": "2307-878787",
  "discarded": "do not return, Please!"
}
`
	mock := mock.Definition{Request: req, Response: res}
	varsProcessor := getProcessor()
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != expected {
		t.Error("Replaced tags from body form do not match", mock.Response.Body)
	}
}

func TestGetBodyParamWithoutContentType(t *testing.T) {
	req := mock.Request{}
	req.Headers = make(mock.Values)
	// No Content-Type header
	req.Body = `
{
  "name": "john",
  "age": 30
}
`
	res := mock.Response{}
	res.Body = `{"name": "{{request.body.name}}", "age": {{request.body.age}}}`

	expected := `{"name": "john", "age": 30}`

	mock := mock.Definition{Request: req, Response: res}
	varsProcessor := getProcessor()
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != expected {
		t.Errorf("Failed to resolve body vars without content-type. Got: %s, Expected: %s", mock.Response.Body, expected)
	}
}

func TestGetBodyParamWithBinaryContentType(t *testing.T) {
	req := mock.Request{}
	req.Headers = make(mock.Values)
	req.Headers["Content-Type"] = []string{"application/octet-stream"}
	req.Body = `{"name": "john"}`
	res := mock.Response{}
	res.Body = `{"name": "{{request.body.name}}"}`

	expected := `{"name": "john"}`

	mock := mock.Definition{Request: req, Response: res}
	varsProcessor := getProcessor()
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != expected {
		t.Errorf("Failed to resolve body vars with application/octet-stream. Got: %s, Expected: %s", mock.Response.Body, expected)
	}
}

func TestURI(t *testing.T) {
	const MOCK_URI = "Test_URI.yml"
	const MOCK_HEADER_NAME = "x-test-uri"

	req := mock.Request{}
	req.Headers = make(mock.Values)
	req.Headers["Content-Type"] = []string{"application/json"}

	res := mock.Response{}
	res.Headers = make(mock.Values)
	res.Headers["Content-Type"] = []string{"application/json"}
	res.Headers[MOCK_HEADER_NAME] = []string{"{{ URI }}"}

	res.Body = `
{
  "URI": "{{ URI }}",
}
`
	var expectedBody = fmt.Sprintf(`
{
  "URI": "%v",
}
`, MOCK_URI)

	mock := mock.Definition{Request: req, Response: res}
	mock.URI = MOCK_URI

	varsProcessor := getProcessor()
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != expectedBody {
		t.Error("failed to replace URI in response body", mock.Response.Body)
	}

	var headerReplaced = slices.Contains(mock.Response.Headers[MOCK_HEADER_NAME], MOCK_URI)
	if !headerReplaced {
		t.Error("failed to replace URI in response headers", mock.Response.Body)
	}
}

func TestDescription(t *testing.T) {
	const MOCK_DESCRIPTION = "TestDescription.yml"
	const MOCK_HEADER_NAME = "x-test-description"

	req := mock.Request{}
	req.Headers = make(mock.Values)
	req.Headers["Content-Type"] = []string{"application/json"}

	res := mock.Response{}
	res.Headers = make(mock.Values)
	res.Headers["Content-Type"] = []string{"application/json"}
	res.Headers[MOCK_HEADER_NAME] = []string{"{{ description }}"}

	res.Body = `
{
  "Description": "{{ description }}",
}
`
	var expectedBody = fmt.Sprintf(`
{
  "Description": "%v",
}
`, MOCK_DESCRIPTION)

	mock := mock.Definition{Request: req, Response: res}
	mock.Description = MOCK_DESCRIPTION

	varsProcessor := getProcessor()
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != expectedBody {
		t.Error("failed to replace Description in response body", mock.Response.Body)
	}

	var headerReplaced = slices.Contains(mock.Response.Headers[MOCK_HEADER_NAME], MOCK_DESCRIPTION)
	if !headerReplaced {
		t.Error("failed to replace Description in response headers", mock.Response.Body)
	}
}
