package vars

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"

	xj "github.com/basgys/goxml2json"
	"github.com/jmartin82/mmock/v3/pkg/mock"
	"github.com/jmartin82/mmock/v3/pkg/route"
	"github.com/tidwall/gjson"
)

type Request struct {
	Mock    *mock.Definition
	Request *mock.Request
}

func (rp Request) Fill(holders []string) map[string][]string {

	vars := make(map[string][]string)
	for _, tag := range holders {
		found := false
		s := ""
		if tag == "request.body" && rp.Request.Body != "" {
			s = rp.Request.Body
			found = true
		} else if tag == "request.scheme" {
			s, found = rp.Request.Scheme, true
		} else if tag == "request.port" {
			s, found = rp.Request.Port, true
		} else if tag == "request.url" {
			s, found = rp.getUrl()
		} else if tag == "request.authority" {
			s, found = rp.getAuthority()
		} else if tag == "request.hostname" {
			s, found = rp.Request.Host, true
		} else if tag == "request.path" {
			s, found = rp.Request.Path, true
		} else if tag == "request.fragment" {
			s, found = rp.Request.Fragment, true
		} else if strings.HasPrefix(tag, "request.body.") {
			s, found = rp.getBodyParam(tag[13:])
		} else if strings.HasPrefix(tag, "request.query.") {
			s, found = rp.getQueryStringParam(tag[14:])
		} else if strings.HasPrefix(tag, "request.path.") {
			s, found = rp.getPathParam(tag[13:])
		} else if strings.HasPrefix(tag, "request.cookie.") {
			s, found = rp.getCookieParam(tag[15:])
		} else if strings.HasPrefix(tag, "request.header.") {
			s, found = rp.getHeaderParam(tag[15:])
		} else if strings.HasPrefix(tag, "env.") {
			s, found = os.LookupEnv(tag[4:])
		} else if strings.ToUpper(tag) == "URI" {
			s, found = rp.Mock.URI, true
		} else if strings.ToLower(tag) == "description" {
			s, found = rp.Mock.Description, true
		}

		if found {
			vars[tag] = append(vars[tag], s)
		}

	}
	return vars
}

func (rp Request) getAuthority() (string, bool) {
	if len(rp.Request.Port) == 0 || rp.Request.Port == "80" {
		return fmt.Sprintf("%s://%s", rp.Request.Scheme, rp.Request.Host), true
	}

	return fmt.Sprintf("%s://%s:%s", rp.Request.Scheme, rp.Request.Host, rp.Request.Port), true
}

func (rp Request) getUrl() (string, bool) {
	value, f := rp.getAuthority()

	if !f {
		return "", false
	}

	path := rp.Request.Path

	if len(path) != 0 {
		value += rp.Request.Path
	}

	queryStringParams := rp.Request.QueryStringParameters

	if len(queryStringParams) != 0 {
		queryKeys := []string{}
		queryVars := []string{}

		//make predictable
		for key := range queryStringParams {
			queryKeys = append(queryKeys, key)
		}
		sort.Strings(queryKeys)

		for _, key := range queryKeys {
			for _, value := range queryStringParams[key] {
				queryVars = append(queryVars, fmt.Sprintf("%s=%s", key, value))
			}
		}
		value += "?" + strings.Join(queryVars, "&")
	}

	if len(rp.Request.Fragment) != 0 {
		value += "#" + rp.Request.Fragment
	}

	return value, true
}

func (rp Request) getPathParam(name string) (string, bool) {

	route := route.NewRoute(rp.Mock.Request.Path)
	mparm := route.Match(rp.Request.Path)

	value, f := mparm.Params[name]
	if !f {
		return "", false
	}

	return value, true
}

func (rp Request) getQueryStringParam(name string) (string, bool) {

	if len(rp.Request.QueryStringParameters) == 0 {
		return "", false
	}
	value, f := rp.Request.QueryStringParameters[name]
	if !f {
		return "", false
	}

	return value[0], true
}

func (rp Request) getCookieParam(name string) (string, bool) {

	if len(rp.Request.Cookies) == 0 {
		return "", false
	}
	value, f := rp.Request.Cookies[name]
	if !f {
		return "", false
	}

	return value, true
}

func (rp Request) getHeaderParam(name string) (string, bool) {

	value, f := rp.Request.HttpHeaders.Headers[name]
	if !f || len(rp.Request.HttpHeaders.Headers) == 0 {
		return "", false
	}
	if len(value) == 0 {
		return "", false
	}

	return value[0], true
}

func (rp Request) getBodyParam(name string) (string, bool) {
	contentType, found := rp.Request.Headers["Content-Type"]
	useSniffing := !found || len(contentType) == 0

	if found && len(contentType) > 0 {
		ct := contentType[0]
		if strings.HasPrefix(ct, "application/x-www-form-urlencoded") {
			return rp.getUrlEncodedFormBodyParam(name)
		} else if strings.HasPrefix(ct, "application/xml") || strings.HasPrefix(ct, "text/xml") {
			return rp.getXmlBodyParam(rp.Request.Body, name)
		} else if strings.HasPrefix(ct, "application/") && strings.HasSuffix(ct, "json") {
			return rp.getJsonBodyParam(rp.Request.Body, name)
		} else if strings.HasPrefix(ct, "application/octet-stream") || strings.HasPrefix(ct, "application/binary") {
			useSniffing = true
		}
	}

	if useSniffing {
		// Content sniffing
		trimmedBody := strings.TrimLeft(rp.Request.Body, " \t\r\n")
		if len(trimmedBody) > 0 {
			firstChar := trimmedBody[0]
			if firstChar == '{' || firstChar == '[' {
				return rp.getJsonBodyParam(rp.Request.Body, name)
			} else if firstChar == '<' {
				return rp.getXmlBodyParam(rp.Request.Body, name)
			}
		}
	}

	return "", false
}

func (rp Request) getXmlBodyParam(body string, name string) (string, bool) {
	xml := strings.NewReader(body)
	json, err := xj.Convert(xml)

	if err != nil {
		return "", false
	}

	value, ret := rp.getBodyParamValue(json.String(), name)

	//TODO: Add support to complex types extraction like arrays or maps
	if value.Type == gjson.JSON {
		return "", false
	}

	return value.String(), ret
}

func (rp Request) getJsonBodyParam(body string, name string) (string, bool) {
	value, ret := rp.getBodyParamValue(body, name)
	return value.String(), ret
}

func (rp Request) getUrlEncodedFormBodyParam(name string) (string, bool) {

	values, err := url.ParseQuery(rp.Request.Body)
	if err != nil {
		return "", false
	}

	value := values.Get(name)
	if value == "" {
		return "", false
	}

	return value, true
}

func (rp Request) haveRegex(query string) bool {
	match, _ := regexp.MatchString("(.regex\\(.*\\))", query)
	return match
}

func (rp Request) haveConcat(query string) bool {
	match, _ := regexp.MatchString("(.concat\\(.*\\))", query)
	return match
}

func (rp Request) getBodyParamValue(body string, query string) (value gjson.Result, found bool) {

	if rp.haveRegex(query) {
		value, found = rp.getParamWithRegex(body, query)
	} else if rp.haveConcat(query) {
		value, found = rp.concatValue(body, query)
	} else {
		value, found = rp.getParamWithGJson(body, query)
	}

	return
}

func (rp Request) getParamWithRegex(body string, query string) (value gjson.Result, found bool) {

	queries := strings.Split(query, ".regex")

	if len(queries) > 1 {

		value, found = rp.getParamWithGJson(body, queries[0])

		regex := queries[1]
		concatValue := ""

		if rp.haveConcat(regex) {
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

func (rp Request) concatValue(body string, query string) (value gjson.Result, found bool) {

	queries := strings.Split(query, ".concat")

	if len(queries) > 1 {

		value, found = rp.getParamWithGJson(body, queries[0])

		concatValue := queries[1]

		if value.Str != "" {
			value.Str += concatValue[1 : len(concatValue)-1]
		}
	}

	return
}

func (rp Request) getParamWithGJson(body string, query string) (value gjson.Result, found bool) {

	value = gjson.Get(body, query)
	if !value.Exists() {
		value.Str = ""
		return value, false
	}

	return value, true
}
