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
	ErrSchemeNotMatch   = errors.New("Scheme not match")
	ErrFragmentNotMatch = errors.New("Fragment not match")
	ErrMethodNotMatch   = errors.New("Method not match")
	ErrPathNotMatch     = errors.New("Path not match")
	ErrQueryStringMatch = errors.New("Query string not match")
	ErrHeadersNotMatch  = errors.New("Headers not match")
	ErrCookiesNotMatch  = errors.New("Cookies not match")
	ErrBodyNotMatch     = errors.New("Body not match")
	ErrScenarioNotMatch = errors.New("Scenario state not match")
)

func NewTester(scenario scenario.Director) *Tester {
	return &Tester{Scenario: scenario}
}

type Tester struct {
	Scenario scenario.Director
}

func (mm Tester) matchKeyAndValues(reqMap definition.Values, mockMap definition.Values) bool {

	if len(mockMap) > len(reqMap) {

		return false
	}

	for key, mval := range mockMap {
		if rval, exists := reqMap[key]; exists {

			if len(mval) > len(rval) {
				return false
			}

			for i, v := range mval {
				if (!strings.Contains(v, glob.GLOB) && v != rval[i]) || !glob.Glob(v, rval[i]) {
					return false
				}
			}

		} else {
			if rval, exists = mm.findByPartialKey(reqMap, key); exists {

				for i, v := range mval {
					if (!strings.Contains(v, glob.GLOB) && v != rval[i]) || !glob.Glob(v, rval[i]) {
						return false
					}
				}
			} else {
				return false
			}
		}
	}
	return true
}

func (mm Tester) findByPartialKey(reqMap definition.Values, partialMatch string) ([]string, bool) {
	if !strings.Contains(partialMatch, glob.GLOB) {
		return []string{}, false
	}

	for key, _ := range reqMap {
		if glob.Glob(partialMatch, key) {
			return reqMap[key], true
		}
	}

	return []string{}, false
}

func (mm Tester) matchKeyAndValue(reqMap definition.Cookies, mockMap definition.Cookies) bool {
	if len(mockMap) > len(reqMap) {
		return false
	}
	for key, mval := range mockMap {
		if rval, exists := reqMap[key]; !exists || (!strings.Contains(mval, glob.GLOB) && mval != rval) || !glob.Glob(mval, rval) {
			return false
		}
	}
	return true
}

func (mm Tester) matchOnEqualsOrIfEmpty(reqVal string, mockVal string) bool {
	if len(mockVal) == 0 {
		return true
	}
	return strings.ToLower(mockVal) == strings.ToLower(reqVal)
}

func (mm Tester) matchOnEqualsOrIfEmptyOrGlob(reqVal string, mockVal string) bool {
	if len(mockVal) == 0 {
		return true
	}
	mockHost := strings.ToLower(mockVal)
	reqHost := strings.ToLower(reqVal)

	return (mockHost == reqHost) || glob.Glob(mockHost, reqHost)
}

func (mm Tester) mockIncludesMethod(method string, mock *definition.Request) bool {
	for _, item := range strings.Split(mock.Method, "|") {
		if strings.ToLower(item) == strings.ToLower(method) {
			return true
		}
	}
	return false
}

func (mm Tester) matchScenarioState(scenario *definition.Scenario) bool {
	if scenario.Name == "" {
		return true
	}

	currentState := mm.Scenario.GetState(scenario.Name)
	for _, r := range scenario.RequiredState {
		if strings.ToLower(r) == currentState {
			return true
		}
	}

	return false
}

func (mm Tester) Check(req *definition.Request, mock *definition.Mock, scenarioAware bool) (bool, error) {

	routes := urlmatcher.New(mock.Request.Path)

	if !mm.matchOnEqualsOrIfEmptyOrGlob(req.Host, mock.Request.Host) {
		return false, ErrHostNotMatch
	}

	if !mm.matchOnEqualsOrIfEmpty(req.Scheme, mock.Request.Scheme) {
		return false, ErrSchemeNotMatch
	}

	if !mm.matchOnEqualsOrIfEmpty(req.Fragment, mock.Request.Fragment) {
		return false, ErrFragmentNotMatch
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

	if scenarioAware && !mm.matchScenarioState(&mock.Control.Scenario) {
		return false, ErrScenarioNotMatch
	}

	return true, nil
}
