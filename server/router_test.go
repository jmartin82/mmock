package server

import (
	"errors"
	"testing"

	"github.com/jmartin82/mmock/definition"
)

type DummyMatcher struct {
	OK bool
}

func (dm DummyMatcher) Check(req *definition.Request, mock *definition.Mock, scenarioAware bool) (bool, error) {
	if dm.OK {
		return true, nil
	}
	return false, errors.New("Random Error")
}

type DummyMapper struct {
	mocks []definition.Mock
}

func (mm DummyMapper) Set(URI string, mock definition.Mock) error {
	return nil
}
func (mm DummyMapper) Delete(URI string) error {
	return nil
}
func (mm DummyMapper) Get(URI string) (definition.Mock, bool) {
	return definition.Mock{}, false
}
func (mm DummyMapper) List() []definition.Mock {
	return mm.mocks
}

func TestValidRoute(t *testing.T) {

	mocks := []definition.Mock{
		{
			Response: definition.Response{
				StatusCode: 200,
			},
		},
	}

	dummyMapper := DummyMapper{mocks: mocks}
	dummyMatcher := DummyMatcher{OK: true}

	r := NewRouter(dummyMapper, dummyMatcher)
	req := definition.Request{Path: "/test"}

	m, result := r.Resolve(&req)

	if len(result.Errors) > 0 || m.Response.StatusCode != 200 {
		t.Fatalf("Not route resolved")
	}

}

func TestInvalidRoute(t *testing.T) {

	mocks := []definition.Mock{
		{
			URI: "XX",
			Response: definition.Response{
				StatusCode: 200,
			},
		},
	}
	dummyMapper := DummyMapper{mocks: mocks}
	dummyMatcher := DummyMatcher{OK: false}

	r := NewRouter(dummyMapper, dummyMatcher)

	req := definition.Request{Path: "/test"}

	_, result := r.Resolve(&req)

	if len(result.Errors) == 0 || result.Errors[0].URI != "XX" || result.Errors[0].Reason != "Random Error" {
		t.Fatalf("Invalid route resolved")
	}

}
