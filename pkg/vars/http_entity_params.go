package vars

import (
        "net/url"
        "regexp"
        "strings"

    "github.com/tidwall/gjson"
    xj "github.com/basgys/goxml2json"

    "github.com/jmartin82/mmock/v3/pkg/mock"
  )

type HttpEntityParams struct {
  Entity  *mock.HTTPEntity
}

func (ep HttpEntityParams) getCookieParam(name string) (string, bool) {

	if len(ep.Entity.HttpHeaders.Cookies) == 0 {
		return "", false
	}
	value, f := ep.Entity.HttpHeaders.Cookies[name]
	if !f {
		return "", false
	}

	return value, true
}

func (ep HttpEntityParams) getHeaderParam(name string) (string, bool) {

	value, f := ep.Entity.HttpHeaders.Headers[name]
	if !f || len(ep.Entity.HttpHeaders.Headers) == 0 {
		return "", false
	}
	if len(value) == 0 {
		return "", false
	}

	return value[0], true
}

func (ep HttpEntityParams) getBodyParam(name string) (string, bool) {
	contentType, found := ep.Entity.Headers["Content-Type"]
	if !found {
		return "", false
	}

	if strings.HasPrefix(contentType[0], "application/x-www-form-urlencoded") {

		return ep.getUrlEncodedFormBodyParam(ep.Entity.Body, name)

	} else if strings.HasPrefix(contentType[0], "application/xml") ||
	  strings.HasPrefix(contentType[0], "text/xml") {

		return ep.getXmlBodyParam(ep.Entity.Body, name)

	} else if strings.HasPrefix(contentType[0], "application/") &&
	  strings.HasSuffix(contentType[0], "json") {

		return ep.getJsonBodyParam(ep.Entity.Body, name)
	}

	return "", false
}


func (ep HttpEntityParams) getXmlBodyParam(body string, name string) (string, bool) {
	xml := strings.NewReader(body)
	json, err := xj.Convert(xml)

	if err != nil {
		return "", false
	}

	value, ret := ep.getBodyParamValue(json.String(), name)

	//TODO: Add support to complex types extraction like arrays or maps
	if value.Type == gjson.JSON {
		return "", false
	}

	return value.String(), ret
}

func (ep HttpEntityParams) getJsonBodyParam(body string, name string) (string, bool) {
	value, ret := ep.getBodyParamValue(body, name)
	return value.String(), ret
}

func (ep HttpEntityParams) getUrlEncodedFormBodyParam(body string, name string) (string, bool) {

	values, err := url.ParseQuery(body)
	if err != nil {
		return "", false
	}

	value := values.Get(name)
	if value == "" {
		return "", false
	}

	return value, true
}

func (ep HttpEntityParams) haveRegex(query string) bool {
	match, _ := regexp.MatchString("(.regex\\(.*\\))", query)
	return match
}

func (ep HttpEntityParams) haveConcat(query string) bool {
	match, _ := regexp.MatchString("(.concat\\(.*\\))", query)
	return match
}

func (ep HttpEntityParams) getBodyParamValue(body string, query string) (value gjson.Result, found bool) {

	if ep.haveRegex(query) {
		value, found = ep.getParamWithRegex(body, query)
	} else if ep.haveConcat(query) {
		value, found = ep.concatValue(body, query)
	} else {
		value, found = ep.getParamWithGJson(body, query)
	}

	return
}

func (ep HttpEntityParams) getParamWithRegex(body string, query string) (value gjson.Result, found bool) {

	queries := strings.Split(query, ".regex")

	if len(queries) > 1 {

		value, found = ep.getParamWithGJson(body, queries[0])

		regex := queries[1]
		concatValue := ""

		if ep.haveConcat(regex) {
			queriesConcat := strings.Split(regex, ".concat")
			regex = queriesConcat[0]
			concatValue = queriesConcat[1]
			concatValue = concatValue[1 : len(concatValue)-1]
		}

		regexValue, err := regexp.Compile(regex[1 : len(regex)-1])
		if err != nil {
			value.Str = ""
			return value, false
		}

		value.Str = regexValue.FindString(value.String())

		if value.Str != "" {
			value.Str += concatValue
		}

	}

	return
}

func (ep HttpEntityParams) concatValue(body string, query string) (value gjson.Result, found bool) {

	queries := strings.Split(query, ".concat")

	if len(queries) > 1 {

		value, found = ep.getParamWithGJson(body, queries[0])

		concatValue := queries[1]

		if value.Str != "" {
			value.Str += concatValue[1 : len(concatValue)-1]
		}
	}

	return
}

func (ep HttpEntityParams) getParamWithGJson(body string, query string) (value gjson.Result, found bool) {

	value = gjson.Get(body, query)
	if !value.Exists() {
		value.Str = ""
		return value, false
	}

	return value, true
}
