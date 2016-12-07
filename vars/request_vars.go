package vars

import (
	"regexp"
	"strings"

	urlmatcher "github.com/azer/url-router"
	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/utils"
)

type RequestVars struct {
	Request     *definition.Request
	RegexHelper utils.RegexHelper
}

func (rp RequestVars) Fill(m *definition.Mock, input string) string {
	r := regexp.MustCompile(`\{\{\s*request\.(.+?)\s*\}\}`)

	return r.ReplaceAllStringFunc(input, func(raw string) string {
		// replace the strings
		if r, found := rp.replaceString(m, raw); found {
			return r
		}
		// replace regexes
		return raw
	})

}

func (rp RequestVars) replaceString(m *definition.Mock, raw string) (string, bool) {
	found := false
	s := ""
	tag := strings.Trim(raw[2:len(raw)-2], " ")
	if tag == "request.body" {
		s = rp.Request.Body
		found = true
	} else if i := strings.Index(tag, "request.query."); i == 0 {
		s, found = rp.getQueryStringParam(rp.Request, tag[14:])
	} else if i := strings.Index(tag, "request.path."); i == 0 {
		s, found = rp.getPathParm(m, rp.Request, tag[13:])
	} else if i := strings.Index(tag, "request.cookie."); i == 0 {
		s, found = rp.getCookieParam(rp.Request, tag[15:])
	}
	if !found {
		return raw, false
	}
	return s, true
}

func (rp RequestVars) getPathParm(m *definition.Mock, req *definition.Request, name string) (string, bool) {

	routes := urlmatcher.New(m.Request.Path)
	mparm := routes.Match(req.Path)

	value, f := mparm.Params[name]
	if !f {
		return "", false
	}

	return value, true
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
