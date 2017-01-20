package match

import (
	"errors"
	"strings"

	"github.com/jmartin82/mmock/scenario"

	urlmatcher "github.com/azer/url-router"
	"github.com/jmartin82/mmock/definition"
	"github.com/ryanuber/go-glob"
)

var (
	ErrHostNotMatch     = errors.New("Host not match")
	ErrMethodNotMatch   = errors.New("Method not match")
	ErrPathNotMatch     = errors.New("Path not match")
	ErrQueryStringMatch = errors.New("Query string not match")
	ErrHeadersNotMatch  = errors.New("Headers not match")
	ErrCookiesNotMatch  = errors.New("Cookies not match")
	ErrBodyNotMatch     = errors.New("Body not match")
	ErrScenarioNotMatch = errors.New("Scenario state not match")
)

type MockMatch struct {
	Scenario scenario.ScenarioManager
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

func (mm MockMatch) mockMatchHost(host string, mock *definition.Request) bool {
	if len(mock.Host) == 0 {
		return true
	}
	return strings.ToLower(host) == strings.ToLower(mock.Host)
}

func (mm MockMatch) mockIncludesMethod(method string, mock *definition.Request) bool {
	for _, item := range strings.Split(mock.Method, "|") {
		if strings.ToLower(item) == strings.ToLower(method) {
			return true
		}
	}
	return false
}

func (mm MockMatch) matchScenarioState(scenario *definition.Scenario) bool {
	if scenario.Name == "" {
		return true
	}

	currentState := mm.Scenario.GetState(scenario.Name)
	for _, r := range scenario.RequiredState {
		if r == currentState {
			return true
		}
	}

	return false
}

func (mm MockMatch) Match(req *definition.Request, mock *definition.Mock) (bool, error) {

	routes := urlmatcher.New(mock.Request.Path)

	if !mm.mockMatchHost(req.Host, &mock.Request) {
		return false, ErrHostNotMatch
	}

	if !glob.Glob(mock.Request.Path, req.Path) && routes.Match(req.Path) == nil {
		return false, ErrPathNotMatch
	}

	if !mm.mockIncludesMethod(req.Method, &mock.Request) {
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
		return false, ErrScenarioNotMatch
	}

	return true, nil
}
