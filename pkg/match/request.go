package match

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/jmartin82/mmock/v3/internal/config/logger"
	"github.com/jmartin82/mmock/v3/pkg/match/payload"
	"github.com/jmartin82/mmock/v3/pkg/mock"
	"github.com/jmartin82/mmock/v3/pkg/route"
	"github.com/ryanuber/go-glob"
)

var log = logger.Log

var (
	ErrCookiesNotMatch  = errors.New("Cookies not match")
	ErrScenarioNotMatch = errors.New("Scenario state not match")
	ErrPathNotMatch     = errors.New("Path not match")
)

func NewTester(comparator *payload.Comparator, scenario ScenearioStorer) *Request {
	return &Request{scenario: scenario, comparator: comparator}
}

type Request struct {
	scenario   ScenearioStorer
	comparator *payload.Comparator
}

func (mm Request) matchKeyAndValues(reqMap mock.Values, mockMap mock.Values) bool {
	if len(mockMap) > len(reqMap) {
		log.Debugf("mock contains more values [%d] than request [%d]",
			len(mockMap), len(reqMap))

		return false
	}

	for key, mval := range mockMap {
		if rval, exists := reqMap[key]; exists {

			if len(mval) > len(rval) {
				log.Debugf("length of mock value [%d] > request value [%d]",
					len(mval), len(rval))

				return false
			}

			if !((mm.matchKey(rval, mval, globMatch)) ||
				(mm.matchKey(rval, mval, regexpMatch))) {

				return false
			}

		} else {

			if rval, exists = mm.findByPartialKey(reqMap, key); exists {

				if !((mm.matchKey(rval, mval, globMatch)) ||
					(mm.matchKey(rval, mval, regexpMatch))) {

					return false
				}
			} else {
				log.Debugf("value %v doesn't appear in mock", key)

				return false
			}
		}
	}
	return true
}

type valueMatcher func(string, string) bool

func globMatch(m string, v string) bool {

	matched := ((strings.Contains(m, glob.GLOB) && glob.Glob(m, v)) || (m == v))
	log.Debugf("value %v globMatch %v: %v", v, m, matched)

	return matched
}

func regexpMatch(m string, v string) bool {
	matched, err := regexp.MatchString(m, fmt.Sprint(v))
	log.Debugf("value %v regexpMatch %v: %v [%v]", v, m, matched, err)

	return (err == nil && matched)
}

func (mm Request) matchKey(rval []string, mval []string, findMatch valueMatcher) bool {
	for i, m := range mval {
		if findMatch(m, rval[i]) {

			return true
		}
	}
	return false
}

func (mm Request) findByPartialKey(reqMap mock.Values, partialMatch string) ([]string, bool) {
	if !strings.Contains(partialMatch, glob.GLOB) {
		return []string{}, false
	}

	for key := range reqMap {
		if glob.Glob(partialMatch, key) {
			return reqMap[key], true
		}
	}

	return []string{}, false
}

func (mm Request) matchKeyAndValue(reqMap mock.Cookies, mockMap mock.Cookies) bool {
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

func (mm Request) matchOnEqualsOrIfEmpty(reqVal string, mockVal string) bool {
	if len(mockVal) == 0 {
		return true
	}
	return strings.EqualFold(mockVal, reqVal)
}

func (mm Request) matchOnEqualsOrIfEmptyOrGlob(reqVal string, mockVal string) bool {
	if len(mockVal) == 0 {
		return true
	}
	mockHost := strings.ToLower(mockVal)
	reqHost := strings.ToLower(reqVal)

	return (mockHost == reqHost) || glob.Glob(mockHost, reqHost)
}

func (mm Request) mockIncludesMethod(method string, mock *mock.Request) bool {
	for _, item := range strings.Split(mock.Method, "|") {
		if strings.EqualFold(item, method) {
			return true
		}
	}

	return false
}

func (mm Request) matchScenarioState(scenario *mock.Scenario) bool {
	if scenario.Name == "" {
		return true
	}

	currentState := mm.scenario.GetState(scenario.Name)
	for _, r := range scenario.RequiredState {
		if strings.ToLower(r) == currentState {
			return true
		}
	}

	return false
}

// Matcher checks if the received request matches with some specific mock request config.
type Matcher interface {
	Match(req *mock.Request, mock *mock.Definition, scenarioAware bool) (bool, error)
}

func (mm Request) Match(req *mock.Request, mock *mock.Definition, scenarioAware bool) (bool, error) {

	route := route.NewRoute(mock.Request.Path)

	if !mm.matchOnEqualsOrIfEmptyOrGlob(req.Host, mock.Request.Host) {
		return false, fmt.Errorf("Host not match. Actual: %s, Expected: %s", req.Host, mock.Request.Host)
	}

	if !mm.matchOnEqualsOrIfEmpty(req.Scheme, mock.Request.Scheme) {
		return false, fmt.Errorf("Scheme not match. Actual: %s, Expected: %s", req.Scheme, mock.Request.Scheme)
	}

	if !mm.matchOnEqualsOrIfEmpty(req.Fragment, mock.Request.Fragment) {
		return false, fmt.Errorf("Fragment not match. Actual: %s, Expected: %s", req.Fragment, mock.Request.Fragment)
	}

	if !glob.Glob(mock.Request.Path, req.Path) && route.Match(req.Path) == nil {
		return false, fmt.Errorf("%w Actual: %s, Expected: %s", ErrPathNotMatch, req.Path, mock.Request.Path)
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

	if !mm.bodyMatch(mock.Request, req) {
		return false, fmt.Errorf("Body not match. Actual: %s, Expected: %s", req.Body, mock.Request.Body)
	}

	if scenarioAware && !mm.matchScenarioState(&mock.Control.Scenario) {
		return false, ErrScenarioNotMatch
	}

	return true, nil
}
func (mm Request) bodyMatch(mockReq mock.Request, req *mock.Request) bool {

	if len(mockReq.Body) == 0 {
		return true
	}

	if mockReq.Body == req.Body {
		return true
	}

	if strings.Contains(mockReq.Body, glob.GLOB) && glob.Glob(mockReq.Body, req.Body) {
		return true
	}

	// Check if we should use a specific comparator based on Content-Type
	contentType, hasContentType := req.Headers["Content-Type"]
	useSniffing := !hasContentType || len(contentType) == 0

	if !useSniffing {
		ct := contentType[0]
		// If it's a known non-sniffable type, try comparing
		if comparable, ok := mm.comparator.Compare(ct, mockReq.Body, req.Body); comparable {
			return ok
		}
		// If it's a generic binary type, allow sniffing
		if strings.HasPrefix(ct, "application/octet-stream") || strings.HasPrefix(ct, "application/binary") {
			useSniffing = true
		}
	}

	if useSniffing {
		// Content sniffing
		trimmedBody := strings.TrimLeft(req.Body, " \t\r\n")
		if len(trimmedBody) > 0 {
			firstChar := trimmedBody[0]
			if firstChar == '{' || firstChar == '[' {
				if comparable, ok := mm.comparator.Compare("application/json", mockReq.Body, trimmedBody); comparable {
					return ok
				}
			} else if firstChar == '<' {
				if comparable, ok := mm.comparator.Compare("application/xml", mockReq.Body, trimmedBody); comparable {
					return ok
				}
			}
		}
	}

	return false
}

func (mm Request) ValuesToString(values mock.Values) string {
	var valuesStr []string

	for name, value := range values {
		name = strings.ToLower(name)
		for _, h := range value {
			valuesStr = append(valuesStr, fmt.Sprintf("%v: %v", name, h))
		}
	}

	return strings.Join(valuesStr, ", ")
}
