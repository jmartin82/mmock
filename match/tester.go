package match

import (
	"errors"
	"strings"
    "fmt"
	"github.com/jmartin82/mmock/scenario"
	urlmatcher "github.com/azer/url-router"
	"github.com/jmartin82/mmock/definition"
	"github.com/ryanuber/go-glob"
)

var (
	ErrCookiesNotMatch  = errors.New("Cookies not match")
	ErrScenarioNotMatch = errors.New("Scenario state not match")
	ErrPathNotMatch     = errors.New("Path not match")
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
		return false, fmt.Errorf("Host not match. Actual: %s, Expected: %s", req.Host, mock.Request.Host)
	}

	if !mm.matchOnEqualsOrIfEmpty(req.Scheme, mock.Request.Scheme) {
		return false, fmt.Errorf("Scheme not match. Actual: %s, Expected: %s", req.Scheme, mock.Request.Scheme)
	}

	if !mm.matchOnEqualsOrIfEmpty(req.Fragment, mock.Request.Fragment) {
		return false, fmt.Errorf("Fragment not match. Actual: %s, Expected: %s", req.Fragment, mock.Request.Fragment)
	}

	if !glob.Glob(mock.Request.Path, req.Path) && routes.Match(req.Path) == nil {
		return false, fmt.Errorf("Path not match. Actual: %s, Expected: %s", req.Path, mock.Request.Path)
	}

	if !mm.mockIncludesMethod(req.Method, &mock.Request) {
		return false, fmt.Errorf("Method not match. Actual: %s, Expected: %s", req.Method, mock.Request.Method)
	}

	if !mm.matchKeyAndValues(req.QueryStringParameters, mock.Request.QueryStringParameters) {
	    return false, fmt.Errorf("Query string not match. Actual: %s, Expected: %s", mm.ValuesToString(req.QueryStringParameters), mm.ValuesToString(mock.Request.QueryStringParameters))
	}

	if !mm.matchKeyAndValue(req.Cookies, mock.Request.Cookies) {
		return false, ErrCookiesNotMatch
	}

	if !mm.matchKeyAndValues(req.Headers, mock.Request.Headers) {
		return false, fmt.Errorf("Headers not match. Actual: %s, Expected: %s", req.Headers, mock.Request.Headers)
	}

	if len(mock.Request.Body) > 0 && !glob.Glob(mock.Request.Body, req.Body) {
		return false, fmt.Errorf("Body not match. Actual: %s, Expected: %s", req.Body, mock.Request.Body)
	}

	if scenarioAware && !mm.matchScenarioState(&mock.Control.Scenario) {
		return false, ErrScenarioNotMatch
	}

	return true, nil
}


func (mm Tester) ValuesToString(values definition.Values) (string) {
    var valuesStr [] string

    for name, value := range values {
       name = strings.ToLower(name)
       for _, h := range value {
         valuesStr = append(valuesStr, fmt.Sprintf("%v: %v", name, h))
       }
    }

    return strings.Join(valuesStr, ", ");
}