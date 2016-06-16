package route

import (
	"errors"
	"log"
	"sync"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/match"
)

var ErrRouteNotFound = errors.New("Mock route not found")

type RequestRouter struct {
	Mocks    []definition.Mock
	Matcher  match.Matcher
	DUpdates chan []definition.Mock
	sync.Mutex
}

func (rr *RequestRouter) Route(req *definition.Request) (*definition.Mock, map[string]string) {
	errors := make(map[string]string)
	rr.Lock()
	defer rr.Unlock()
	for _, mock := range rr.Mocks {
		m, err := rr.Matcher.Match(req, &mock.Request)
		if m {
			return &mock, nil
		}
		errors[mock.Name] = err.Error()
		log.Printf("Discarting mock: %s Reason: %s\n", mock.Name, err.Error())

	}

	return &definition.Mock{Response: definition.Response{StatusCode: 404}}, errors

}

func (rr *RequestRouter) SetMockDefinitions(mocks []definition.Mock) {
	rr.Lock()
	rr.Mocks = mocks
	rr.Unlock()
}

func (rr *RequestRouter) MockChangeWatch() {
	go func() {
		for {
			mocks := <-rr.DUpdates
			rr.SetMockDefinitions(mocks)
			log.Println("New mock definitions loaded")
		}

	}()
}
