package vars

import (
	"regexp"
	"strings"

	"github.com/jmartin82/mmock/definition"
)

type RequestVars struct {
	Request *definition.Request
}

func (rp RequestVars) Fill(m *definition.Mock, input string) string {
	r := regexp.MustCompile(`\{\{\s*([^\}]+)\s*\}\}`)

	return r.ReplaceAllStringFunc(input, func(raw string) string {
		found := false
		s := ""
		tag := strings.Trim(raw[2:len(raw)-2], " ")
		if tag == "request.body" {
			s = rp.Request.Body
			found = true
		} else if i := strings.Index(tag, "request.url."); i == 0 {
			s, found = rp.getStringPart(rp.Request.Path, tag[12:], "value")
		} else if i := strings.Index(tag, "request.body."); i == 0 {
			s, found = rp.getStringPart(rp.Request.Body, tag[13:], "value")
		} else if i := strings.Index(tag, "request.query."); i == 0 {
			s, found = rp.getQueryStringParam(rp.Request, tag[14:])
		} else if i := strings.Index(tag, "request.cookie."); i == 0 {
			s, found = rp.getCookieParam(rp.Request, tag[15:])
		}
		if !found {
			return raw
		}
		return s
	})
}

func (rp RequestVars) getStringPart(input string, pattern string, groupName string) (string, bool) {
	r, error := regexp.Compile(pattern)
	if error != nil {
		return "", false
	}

	match := r.FindStringSubmatch(input)
	result := make(map[string]string)
	names := r.SubexpNames()
	if len(match) >= len(names) {
		for i, name := range names {
			if i != 0 {
				result[name] = match[i]
			}
		}
	}

	value, present := result[groupName]

	return value, present
}

func (rp RequestVars) getQueryStringParam(req *definition.Request, name string) (string, bool) {

	if len(rp.Request.QueryStringParameters) == 0 {
		return "", false
	}
	value, f := rp.Request.QueryStringParameters[name]
	if !f {
		return "", false
	}

	return value[0], true
}

func (rp RequestVars) getCookieParam(req *definition.Request, name string) (string, bool) {

	if len(rp.Request.Cookies) == 0 {
		return "", false
	}
	value, f := rp.Request.Cookies[name]
	if !f {
		return "", false
	}

	return value, true
}
