package match

import (
	"errors"
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

func (this MockMatch) matchKeyAndValues(reqMap definition.Values, mockMap definition.Values) bool {
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

func (this MockMatch) matchKeyAndValue(reqMap definition.Cookies, mockMap definition.Cookies) bool {
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

func (this MockMatch) Match(req *definition.Request, mock *definition.Request) (bool, error) {

	if req.Method != mock.Method {
		return false, ErrMethodNotMatch
	}

	if !glob.Glob(mock.Path, req.Path) {
		return false, ErrPathNotMatch
	}

	if !this.matchKeyAndValues(req.QueryStringParameters, mock.QueryStringParameters) {
		return false, ErrQueryStringMatch
	}

	if !this.matchKeyAndValue(req.Cookies, mock.Cookies) {
		return false, ErrCookiesNotMatch
	}

	if !this.matchKeyAndValues(req.Headers, mock.Headers) {
		return false, ErrHeadersNotMatch
	}

	if !glob.Glob(mock.Body, req.Body) {
		return false, ErrHeadersNotMatch
	}

	return true, nil
}
