package vars

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	urlmatcher "github.com/azer/url-router"
	"github.com/jmartin82/mmock/definition"
)

type Request struct {
	Mock    *definition.Mock
	Request *definition.Request
}

func (rp Request) Fill(holders []string) map[string]string {

	vars := make(map[string]string)
	for _, tag := range holders {
		found := false
		s := ""
		if tag == "request.body" && rp.Request.Body != "" {
			s = rp.Request.Body
			found = true
		} else if strings.HasPrefix(tag, "request.body.") {
			s, found = rp.getBodyParam(rp.Request, tag[13:])
		} else if strings.HasPrefix(tag, "request.query.") {
			s, found = rp.getQueryStringParam(rp.Request, tag[14:])
		} else if strings.HasPrefix(tag, "request.path.") {
			s, found = rp.getPathParam(rp.Mock, rp.Request, tag[13:])
		} else if strings.HasPrefix(tag, "request.cookie.") {
			s, found = rp.getCookieParam(rp.Request, tag[15:])
		}

		if found {
			vars[tag] = s
		}

	}
	return vars
}

func (rp Request) getPathParam(m *definition.Mock, req *definition.Request, name string) (string, bool) {

	routes := urlmatcher.New(m.Request.Path)
	mparm := routes.Match(req.Path)

	value, f := mparm.Params[name]
	if !f {
		return "", false
	}

	return value, true
}

func (rp Request) getQueryStringParam(req *definition.Request, name string) (string, bool) {

	if len(rp.Request.QueryStringParameters) == 0 {
		return "", false
	}
	value, f := rp.Request.QueryStringParameters[name]
	if !f {
		return "", false
	}

	return value[0], true
}

func (rp Request) getCookieParam(req *definition.Request, name string) (string, bool) {

	if len(rp.Request.Cookies) == 0 {
		return "", false
	}
	value, f := rp.Request.Cookies[name]
	if !f {
		return "", false
	}

	return value, true
}

func (rp Request) getBodyParam(req *definition.Request, name string) (string, bool) {

	contentType, found := req.Headers["Content-Type"]
	if !found {
		return "", false
	}

	if strings.HasPrefix(contentType[0], "application/x-www-form-urlencoded") {
		return rp.getUrlEncodedFormBodyParam(rp.Request, name)
	} else if strings.HasPrefix(contentType[0], "application/json") {
		return rp.getJsonBodyParam(rp.Request, name)
	}

	return "", false
}

func (rp Request) getUrlEncodedFormBodyParam(req *definition.Request, name string) (string, bool) {

	values, err := url.ParseQuery(req.Body)
	if err != nil {
		return "", false
	}

	value := values.Get(name)
	if value == "" {
		return "", false
	}

	return value, true
}

func (rp Request) getJsonBodyParam(req *definition.Request, name string) (string, bool) {

	hierarchy := strings.Split(name, ".")

	var payload interface{}
	if err := json.Unmarshal([]byte(req.Body), &payload); err != nil {
		return "", false
	}

	for _, value := range hierarchy {
		if mapper, ok := payload.(map[string]interface{}); ok {
			payload, ok = mapper[value]

			if !ok {
				return "", false
			}

			continue
		}

		if arrayMapper, ok := payload.([]interface{}); ok {
			index, err := strconv.Atoi(value)

			if err != nil {
				return "", false
			}

			payload = arrayMapper[index]
			continue
		}
	}

	return rp.getJsonValue(payload)
}

func (rp Request) getJsonValue(object interface{}) (string, bool) {

	stringContent, ok := object.(string)

	if ok {
		return stringContent, true
	}

	genericContent, err := json.Marshal(object)

	if err == nil {
		return string(genericContent), true
	}

	return "", false
}
