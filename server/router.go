package server

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/match"
	"strings"
	"sort"
)

//Resolver contains the functions to check the http request and return the matching mock.
type Resolver interface {
	Resolve(req *definition.Request) (*definition.Mock, *definition.MatchResult)
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
func (rr *Router) Resolve(req *definition.Request) (*definition.Mock, *definition.MatchResult) {
	mocks := rr.Mapping.List()
	mLog := &definition.MatchResult{Found: false}
	mLog.Errors = make([]definition.MatchError, 0, len(mocks))

	for _, mock := range OrderMocksByNumberOfWildcardParams(mocks) {
		m, err := rr.Checker.Check(req, &mock, true)
		if m {
			//we return a copy of it, not the definition itself because we will working on it.
			md := definition.Mock{}
			rr.copy(&mock, &md)
			mLog.Found = true
			mLog.URI = mock.URI
			return &md, mLog
		}
		mLog.Errors = append(mLog.Errors, definition.MatchError{URI: mock.URI, Reason: err.Error()})
		if err != match.ErrPathNotMatch {
			log.Printf("Discarding mock: %s Reason: %s\n", mock.URI, err.Error())
		}
	}

	return getNotFoundResult(), mLog
}

func getNotFoundResult() *definition.Mock {
	return &definition.Mock{Response: definition.Response{StatusCode: 404}}
}

//Orders the mocks by number of wildcard params (*) in order to give more priority to the more restrictive ones
func OrderMocksByNumberOfWildcardParams(mocks []definition.Mock) []definition.Mock {

	var mocksOrderedByWildcardParametersNum []definition.Mock

	for key, _ := range mocks {
		mocksOrderedByWildcardParametersNum = append(mocksOrderedByWildcardParametersNum, mocks[key])
	}

	sort.Slice(mocksOrderedByWildcardParametersNum, func(i, j int) bool {
		currMock := mocksOrderedByWildcardParametersNum[i]
		prevMock := mocksOrderedByWildcardParametersNum[j]

		return GetMockWildcardParametersNumber(currMock) <= GetMockWildcardParametersNumber(prevMock)
	})

	return mocksOrderedByWildcardParametersNum
}

func GetMockWildcardParametersNumber(mock definition.Mock) int {
	currMockWildcardsNum := strings.Count(mock.Request.Body, "*")
	for _, queryStringParam := range mock.Request.QueryStringParameters {
		for _, q := range queryStringParam {
			if strings.Count(q, "*") > 0 {
				currMockWildcardsNum += 1
			}
		}
	}

	return currMockWildcardsNum
}
