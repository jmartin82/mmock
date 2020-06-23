package vars

import (
	"github.com/jmartin82/mmock/v3/pkg/mock"
	"os"
	"strings"
	"testing"
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
