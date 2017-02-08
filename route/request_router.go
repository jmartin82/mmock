package route

import (
	"bytes"
	"encoding/gob"
	"log"
	"sync"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/match"
)

//NewRouter returns a pointer to new RequestRouter
func NewRouter(mocks []definition.Mock, checker match.Checker, dUpdates chan []definition.Mock) *RequestRouter {
	return &RequestRouter{
		Mocks:    mocks,
		Checker:  checker,
		DUpdates: dUpdates,
	}
}

//RequestRouter checks http requesta and try to figure out what is the best mock for each one.
type RequestRouter struct {
	Mocks    []definition.Mock
	Checker  match.Checker
	DUpdates chan []definition.Mock
	sync.Mutex
}

func (rr *RequestRouter) copy(src, dst *definition.Mock) {
	var mod bytes.Buffer
	enc := gob.NewEncoder(&mod)
	dec := gob.NewDecoder(&mod)
	err := enc.Encode(src)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	err = dec.Decode(dst)
	if err != nil {
		log.Fatal("decode error:", err)
	}

}

//Route checks the request with all available mock definitions and return the matching mock for it.
func (rr *RequestRouter) Route(req *definition.Request) (*definition.Mock, definition.MatchErrors) {
	errors := make(definition.MatchErrors)
	rr.Lock()
	defer rr.Unlock()
	for _, mock := range rr.Mocks {
		m, err := rr.Checker.Check(req, &mock, true)
		if m {
			//we return a copy of it, not the definition itself because we will working on it.
			md := definition.Mock{}
			rr.copy(&mock, &md)
			return &md, nil
		}
		errors[mock.Name] = err.Error()
		if err != match.ErrPathNotMatch {
			log.Printf("Discarding mock: %s Reason: %s\n", mock.Name, err.Error())
		}
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
