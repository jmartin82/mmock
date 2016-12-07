package route

import (
	"bytes"
	"encoding/gob"
	"sync"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/logging"
	"github.com/jmartin82/mmock/match"
)

//NewRouter returns a pointer to new RequestRouter
func NewRouter(mocks []definition.Mock, matcher match.Matcher, dUpdates chan []definition.Mock) *RequestRouter {
	return &RequestRouter{
		Mocks:    mocks,
		Matcher:  match.MockMatch{},
		DUpdates: dUpdates,
	}
}

//RequestRouter checks http requesta and try to figure out what is the best mock for each one.
type RequestRouter struct {
	Mocks    []definition.Mock
	Matcher  match.Matcher
	DUpdates chan []definition.Mock
	sync.Mutex
}

func (rr *RequestRouter) Copy(src, dst *definition.Mock) {
	var mod bytes.Buffer
	enc := gob.NewEncoder(&mod)
	dec := gob.NewDecoder(&mod)
	err := enc.Encode(src)
	if err != nil {
		logging.Fatal("encode error:", err)
	}
	err = dec.Decode(dst)
	if err != nil {
		logging.Fatal("decode error:", err)
	}

}

//Route checks the request with all available mock definitions and return the matching mock for it.
func (rr *RequestRouter) Route(req *definition.Request) (*definition.Mock, map[string]string) {
	errors := make(map[string]string)
	rr.Lock()
	defer rr.Unlock()
	for _, mock := range rr.Mocks {
		m, err := rr.Matcher.Match(req, &mock.Request)
		if m {
			//we return a copy of it, not the definition itself because we will working on it.
			md := definition.Mock{}
			rr.Copy(&mock, &md)
			return &md, nil
		}
		errors[mock.Name] = err.Error()
		if err != match.ErrPathNotMatch {
			logging.Printf("Discarding mock: %s Reason: %s\n", mock.Name, err.Error())
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
			logging.Println("New mock definitions loaded")
		}

	}()
}
