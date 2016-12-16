package match

import (
	"errors"
	"strings"

	urlmatcher "github.com/azer/url-router"
	"github.com/jmartin82/mmock/definition"
	"github.com/ryanuber/go-glob"
)

var (
	ErrMethodNotMatch   = errors.New("Method not match")
	ErrPathNotMatch     = errors.New("Path not match")
	ErrQueryStringMatch = errors.New("Query string not match")
	ErrHeadersNotMatch  = errors.New("Headers not match")
	ErrCookiesNotMatch  = errors.New("Cookies not match")
	ErrBodyNotMatch     = errors.New("Body not match")
)

type MockMatch struct {
}

func (mm MockMatch) matchKeyAndValues(reqMap definition.Values, mockMap definition.Values) bool {

	if len(mockMap) > len(reqMap) {

		return false
	}

	for key, mval := range mockMap {
		if rval, exists := reqMap[key]; exists {

			if len(mval) > len(rval) {
				return false
			}

			for i, v := range mval {
				if v != rval[i] {
					return false
				}
			}

		} else {
			return false
		}
	}
	return true
}

func (mm MockMatch) matchKeyAndValue(reqMap definition.Cookies, mockMap definition.Cookies) bool {
	if len(mockMap) > len(reqMap) {
		return false
	}
	for key, mval := range mockMap {
		if rval, exists := reqMap[key]; !exists || mval != rval {
			return false
		}
	}
	return true
}

func (mm MockMatch) mockIncludesMethod(mock *definition.Request, method string) bool {
	for _, item := range strings.Split(mock.Method, "|") {
		if item == method {
			return true
		}
	}
	return false
}
func (mm MockMatch) matchScenarioState(scenario *definition.Scenario) bool {

	return false
}

func (mm MockMatch) Match(req *definition.Request, mock *definition.Mock) (bool, error) {

	routes := urlmatcher.New(mock.Request.Path)

	if routes.Match(req.Path) == nil {
		return false, ErrPathNotMatch
	}

	if !mm.mockIncludesMethod(&mock.Request, req.Method) {
		return false, ErrMethodNotMatch
	}

	if !mm.matchKeyAndValues(req.QueryStringParameters, mock.Request.QueryStringParameters) {
		return false, ErrQueryStringMatch
	}

	if !mm.matchKeyAndValue(req.Cookies, mock.Request.Cookies) {
		return false, ErrCookiesNotMatch
	}

	if !mm.matchKeyAndValues(req.Headers, mock.Request.Headers) {
		return false, ErrHeadersNotMatch
	}

	if len(mock.Request.Body) > 0 && !glob.Glob(mock.Request.Body, req.Body) {
		return false, ErrBodyNotMatch
	}

	if !mm.matchScenarioState(&mock.Control.Scenario) {
		return false, ErrBodyNotMatch
	}

	return true, nil
}
