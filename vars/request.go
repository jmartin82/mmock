package vars

import (
	"strings"

	urlmatcher "github.com/azer/url-router"
	"github.com/jmartin82/mmock/definition"
	"net/url"
	"encoding/json"
	"strconv"
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
		} else if i := strings.Index(tag, "request.body."); i == 0 {
			s, found = rp.getBodyParam(rp.Request, tag[13:])
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

func (rp Request) getPathParm(m *definition.Mock, req *definition.Request, name string) (string, bool) {

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

	if i := strings.Index(contentType[0], "application/x-www-form-urlencoded"); i == 0 {
		return rp.getUrlEncodedFormBodyParam(rp.Request, name)
	} else if i := strings.Index(contentType[0], "application/json"); i == 0 {
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

	var payload map[string]*json.RawMessage
	if err := json.Unmarshal([]byte(req.Body), &payload); err != nil {
		return "", false
	}

	return rp.getJsonObjectParamRecursive(payload, name, 100)
}

func (rp Request) getJsonObjectParamRecursive(payload map[string]*json.RawMessage, name string, levelsLeft uint) (string, bool) {

	parts := strings.SplitN(name, ".", 2)
	value, found := payload[parts[0]]; 
	if !found {
		return "", false
	}

	if len(parts) == 1 {
		return rp.getJsonValue(value)
	}

	levelsLeft = levelsLeft - 1
	if levelsLeft == 0 {
		return "", false
	}

	var nested map[string]*json.RawMessage
	if err := json.Unmarshal(*value, &nested); err != nil {
		return "", false
	}

	return rp.getJsonObjectParamRecursive(nested, parts[1], levelsLeft)
}

func (rp Request) getJsonValue(value *json.RawMessage) (string, bool) {
	// Is json 'null' value
	if value == nil {
		return "null", true
	}

	var toNumber json.Number
	if err := json.Unmarshal(*value, &toNumber); err == nil {
		return toNumber.String(), true
	}

	var toBoolean bool
	if err := json.Unmarshal(*value, &toBoolean); err == nil {
		return strconv.FormatBool(toBoolean), true
	}

	var toString string
	if err := json.Unmarshal(*value, &toString); err == nil {
		return toString, true
	}

	var toArray []*json.RawMessage
	if err := json.Unmarshal(*value, &toArray); err == nil {
		return string(*value), true 
	}

	var toObject map[string]*json.RawMessage
	if err := json.Unmarshal(*value, &toObject); err == nil {
		return string(*value), true
	}

	return "", false
}