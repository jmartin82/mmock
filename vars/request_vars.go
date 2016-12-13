package vars

import (
	"strings"

	urlmatcher "github.com/azer/url-router"
	"github.com/jmartin82/mmock/definition"
)

type RequestVars struct {
	Mock    *definition.Mock
	Request *definition.Request
}

func (rp RequestVars) Fill(holders []string) map[string]string {

	vars := make(map[string]string)
	for _, tag := range holders {
		found := false
		s := ""
		if tag == "request.body" {
			s = rp.Request.Body
			found = true
		} else if i := strings.Index(tag, "request.query."); i == 0 {
			s, found = rp.getQueryStringParam(rp.Request, tag[14:])
		} else if i := strings.Index(tag, "request.path."); i == 0 {
			s, found = rp.getPathParm(rp.Mock, rp.Request, tag[13:])
		} else if i := strings.Index(tag, "request.cookie."); i == 0 {
			s, found = rp.getCookieParam(rp.Request, tag[15:])
		}

		if found {
			vars[tag] = s
		}

	}
	return vars
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
