package route

import (
	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/match"
	"log"
	"sync"
)

//RequestRouter checks http requesta and try to figure out what is the best mock for each one.
type RequestRouter struct {
	Mocks    []definition.Mock
	Matcher  match.Matcher
	DUpdates chan []definition.Mock
	sync.Mutex
}

//Route checks the request with all available mock definitions and return the matching mock for it.
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

//SetMockDefinitions allows replace the current mock definitions for new ones.
func (rr *RequestRouter) SetMockDefinitions(mocks []definition.Mock) {
	rr.Lock()
	rr.Mocks = mocks
	rr.Unlock()
}

//MockChangeWatch monitors the mock configuration dir and loads again all the mocks it something change.
func (rr *RequestRouter) MockChangeWatch() {
	go func() {
		for {
			mocks := <-rr.DUpdates
			rr.SetMockDefinitions(mocks)
			log.Println("New mock definitions loaded")
		}

	}()
}
