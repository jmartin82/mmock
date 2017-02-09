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

func TestValidRoute(t *testing.T) {

	mocks := []definition.Mock{
		definition.Mock{
			Response: definition.Response{
				StatusCode: 200,
			},
		},
	}
	dUpdate := make(chan []definition.Mock)
	dummyMatcher := DummyMatcher{OK: true}

	r := NewRouter(mocks, dummyMatcher, dUpdate)

	req := definition.Request{Path: "/test"}

	m, errs := r.Resolve(&req)

	if len(errs) > 0 || m.Response.StatusCode != 200 {
		t.Fatalf("Not route resolved")
	}

}

func TestInvalidRoute(t *testing.T) {

	mocks := []definition.Mock{
		definition.Mock{
			Name: "XX",
			Response: definition.Response{
				StatusCode: 200,
			},
		},
	}
	dUpdate := make(chan []definition.Mock)
	dummyMatcher := DummyMatcher{OK: false}

	r := NewRouter(mocks, dummyMatcher, dUpdate)

	req := definition.Request{Path: "/test"}

	_, errs := r.Resolve(&req)

	if len(errs) == 0 || errs["XX"] != "Random Error" {
		t.Fatalf("Invalid route resolved")
	}

}

func TestRoutesLoadViaChannel(t *testing.T) {

	mocks := []definition.Mock{
		definition.Mock{
			Name: "XX",
			Response: definition.Response{
				StatusCode: 200,
			},
		},
	}
	dUpdate := make(chan []definition.Mock)
	dummyMatcher := DummyMatcher{OK: true}

	r := NewRouter([]definition.Mock{}, dummyMatcher, dUpdate)

	req := definition.Request{Path: "/test"}

	m, _ := r.Resolve(&req)

	if m.Name != "" {
		t.Fatalf("Invalid route resolved")
	}

	go r.MockChangeWatch()
	dUpdate <- mocks

	m, _ = r.Resolve(&req)

	if m.Name == "" {
		t.Fatalf("Route not found")
	}
}
