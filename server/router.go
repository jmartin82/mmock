package server

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/match"
)

//Resolver contains the functions to check the http request and return the matching mock.
type Resolver interface {
	Resolve(req *definition.Request) (*definition.Mock, definition.MatchErrors)
}

//NewRouter returns a pointer to new Router
func NewRouter(mapping definition.Mapping, checker match.Checker) *Router {
	return &Router{
		Mapping: mapping,
		Checker: checker,
	}
}

//Router checks http requesta and try to figure out what is the best mock for each one.
type Router struct {
	Mapping definition.Mapping
	Checker match.Checker
}

func (rr *Router) copy(src, dst *definition.Mock) {
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
func (rr *Router) Resolve(req *definition.Request) (*definition.Mock, definition.MatchErrors) {
	errors := make(definition.MatchErrors)
	mocks := rr.Mapping.List()

	for _, mock := range mocks {
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
	return getNotFoundResult(), errors
}

func getNotFoundResult() *definition.Mock {
	return &definition.Mock{Response: definition.Response{StatusCode: 404}}
}
