package vars

import (
	"regexp"
	"strings"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/utils"
)

type RequestVars struct {
	Request     *definition.Request
	RegexHelper utils.RegexHelper
}

func (rp RequestVars) Fill(m *definition.Mock, input string, multipleMatch bool) string {
	r := regexp.MustCompile(`\{\{\s*request\.(.+?)\s*\}\}`)

	if !multipleMatch {
		return r.ReplaceAllStringFunc(input, func(raw string) string {
			// replace the strings
			if raw, found := rp.replaceString(raw); found {
				return raw
			}
			// replace regexes
			return rp.replaceRegex(raw)
		})
	} else {
		// first replace all strings
		input = r.ReplaceAllStringFunc(input, func(raw string) string {
			item, _ := rp.replaceString(raw)
			return item
		})
		// get multiple entities using regex
		if results, found := rp.RegexHelper.GetCollectionItems(input, rp.getVarsRegexParts); found {
			if len(results) == 1 {
				return "," + results[0] // add a comma in the beginning so that we will now that the item is a single entity
			}

			return strings.Join(results, ",")
		}
		return input
	}
}

func (rp RequestVars) replaceString(raw string) (string, bool) {
	found := false
	s := ""
	tag := strings.Trim(raw[2:len(raw)-2], " ")
	if tag == "request.body" {
		s = rp.Request.Body
		found = true
	} else if i := strings.Index(tag, "request.query."); i == 0 {
		s, found = rp.getQueryStringParam(rp.Request, tag[14:])
	} else if i := strings.Index(tag, "request.cookie."); i == 0 {
		s, found = rp.getCookieParam(rp.Request, tag[15:])
	}
	if !found {
		return raw, false
	}
	return s, true
}

func (rp RequestVars) getVarsRegexParts(input string) (string, string, bool) {
	if i := strings.Index(input, "request.url."); i == 0 {
		return rp.Request.Path, input[12:], true
	} else if i := strings.Index(input, "request.body."); i == 0 {
		return rp.Request.Body, input[13:], true
	}
	return "", "", false
}

func (rp RequestVars) replaceRegex(raw string) string {
	tag := strings.Trim(raw[2:len(raw)-2], " ")
	if regexInput, regexPattern, found := rp.getVarsRegexParts(tag); found {
		if result, found := rp.RegexHelper.GetStringPart(regexInput, regexPattern, "value"); found {
			return result
		}
	}
	return raw
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
