package route

import (
	"errors"
	"log"
	"mmock/definition"
	"mmock/match"
	"sync"
)

var ErrRouteNotFound = errors.New("Mock route not found")

type RequestRouter struct {
	Mocks    []definition.Mock
	Matcher  match.Matcher
	DUpdates chan []definition.Mock
	sync.Mutex
}

func (this *RequestRouter) Route(req *definition.Request) (*definition.Mock, map[string]string) {
	errors := make(map[string]string)
	this.Lock()
	defer this.Unlock()
	for _, mock := range this.Mocks {
		m, err := this.Matcher.Match(req, &mock.Request)
		if m {
			return &mock, nil
		}
		errors[mock.Name] = err.Error()
		log.Printf("Discarting mock: %s Reason: %s\n", mock.Name, err.Error())

	}

	return &definition.Mock{Response: definition.Response{StatusCode: 404}}, errors

}

func (this *RequestRouter) SetMockDefinitions(mocks []definition.Mock) {
	this.Lock()
	this.Mocks = mocks
	this.Unlock()
}

func (this *RequestRouter) MockChangeWatch() {
	go func() {
		for {
			mocks := <-this.DUpdates
			this.SetMockDefinitions(mocks)
			log.Println("New mock definitions loaded")
		}

	}()
}
