package vars

import (
	"fmt"
	"github.com/jmartin82/mmock/v3/internal/config/logger"
	"github.com/jmartin82/mmock/v3/pkg/mock"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var log = logger.Log

var re = regexp.MustCompile(`(?m)\((.*)\)`)

type Stream struct {
	Mock    *mock.Definition
	Request *mock.Request
}

func (st Stream) Fill(holders []string) map[string][]string {

	vars := make(map[string][]string)
	for _, tag := range holders {
		if strings.HasPrefix(tag, "file.contents(") {
			vars[tag] = append(vars[tag], st.getOutput(st.getFileContents(tag)))
		} else if strings.HasPrefix(tag, "http.contents(") {
			vars[tag] = append(vars[tag], st.getOutput(st.getHttpContents(tag)))
		}
	}
	return vars
}

func (st Stream) getOutput(o []byte, err error) string {

	if err != nil {
		log.Errorf("Impossible read mock stream: %s", err)
		return fmt.Sprintf("ERROR: %s", err.Error())
	}
	return string(o)
}

func (st Stream) getFileContents(tag string) ([]byte, error) {
	path := st.getInputParam(tag, true)
	return ioutil.ReadFile(path)
}

func (st Stream) getHttpContents(tag string) ([]byte, error) {
	url := st.getInputParam(tag, false)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (st Stream) getInputParam(param string, forFilePath bool) string {
	match := re.FindStringSubmatch(param)
	if len(match) > 1 {
		return st.resolveRequestVars(match[1], forFilePath)
	}
	return ""
}

// resolveRequestVars replaces {{request.*}} holders inside a stream argument using
// the same extraction logic as the response templating (path params, JSON body via
// gjson, query params, etc.), so a file path can be composed from request values
// (e.g. file.contents(/data/{{request.path.tenant}}/{{request.body.name}}.csv)).
//
// It is a no-op without request context or holders, so literal paths keep working.
// Interpolation is enabled only for file paths: http.contents URLs are left literal
// to avoid introducing an SSRF surface from client-controlled request data.
func (st Stream) resolveRequestVars(input string, forFilePath bool) string {
	if !forFilePath {
		return input
	}
	if st.Mock == nil || st.Request == nil || !strings.Contains(input, "{{") {
		return input
	}

	rp := Request{Mock: st.Mock, Request: st.Request}

	return varsRegex.ReplaceAllStringFunc(input, func(match string) string {
		token := strings.Trim(match, "{} ")
		if !strings.HasPrefix(token, "request.") {
			return match
		}

		values, found := rp.Fill([]string{token})[token]
		if !found || len(values) == 0 {
			return match
		}

		// Reject values that are not a single, safe path segment (path separators
		// or traversal tokens) so a crafted request cannot escape the base directory.
		if value := values[0]; isSafePathSegment(value) {
			return value
		}
		return match
	})
}

// isSafePathSegment reports whether s is a single path segment safe to embed in a
// file path: it contains no path separator and is not a directory traversal token.
func isSafePathSegment(s string) bool {
	return s != "." && s != ".." && !strings.ContainsAny(s, `/\`)
}
