package parse

import (
	"encoding/json"
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/Jeffail/gabs"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/parse/fakedata"
)

//FakeDataParse parses the data looking for fake data tags or request data tags
type FakeDataParse struct {
	Fake fakedata.DataFaker
}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func (fdp FakeDataParse) call(data reflect.Value, name string) string {
	// get a reflect.Value for the method
	methodVal := data.MethodByName(name)
	// turn that into an interface{}
	methodIface := methodVal.Interface()
	// turn that into a function that has the expected signature
	method := methodIface.(func() string)
	// call the method directly
	res := method()
	return res
}

func (fdp FakeDataParse) callMethod(name string) (string, bool) {
	found := false
	data := reflect.ValueOf(fdp.Fake)
	typ := data.Type()
	if nMethod := data.Type().NumMethod(); nMethod > 0 {
		for i := 0; i < nMethod; i++ {
			method := typ.Method(i)
			if strings.ToLower(method.Name) == strings.ToLower(name) {
				found = true // we found the name regardless
				// does receiver type match? (pointerness might be off)
				if typ == method.Type.In(0) {
					return fdp.call(data, method.Name), found
				}
			}
		}
	}
	return "", found
}

//ReplaceVars relplaces variables from the request in the input
func (fdp FakeDataParse) ReplaceVars(req *definition.Request, res *definition.Response, input string) string {
	r := regexp.MustCompile(`\{\{\s*([^\}]+)\s*\}\}`)

	return r.ReplaceAllStringFunc(input, func(raw string) string {
		found := false
		s := ""
		tag := strings.Trim(raw[2:len(raw)-2], " ")
		if tag == "request.body" {
			s = req.Body
			found = true
		} else if tag == "response.body" {
			s = res.Body
			found = true
		} else if i := strings.Index(tag, "request.url."); i == 0 {
			s, found = getStringPart(req.Path, tag[12:], "value")
		} else if i := strings.Index(tag, "request.body."); i == 0 {
			s, found = getStringPart(req.Body, tag[13:], "value")
		} else if i := strings.Index(tag, "response.body."); i == 0 {
			s, found = getStringPart(res.Body, tag[14:], "value")
		} else if i := strings.Index(tag, "request.query."); i == 0 {
			s, found = getQueryStringParam(req, tag[14:])
		} else if i := strings.Index(tag, "request.cookie."); i == 0 {
			s, found = getCookieParam(req, tag[15:])
		} else if i := strings.Index(tag, "fake."); i == 0 {
			s, found = fdp.callMethod(tag[5:])
		}

		if !found {
			log.Printf("Defined tag {{%s}} not found\n", tag)
			return raw
		}
		return s
	})
}

//Parse subtitutes the current mock response and replace the tags stored inside.
func (fdp FakeDataParse) Parse(req *definition.Request, res *definition.Response) {
	for header, values := range res.Headers {
		for i, value := range values {
			res.Headers[header][i] = fdp.ReplaceVars(req, res, value)
		}

	}
	for cookie, value := range res.Cookies {
		res.Cookies[cookie] = fdp.ReplaceVars(req, res, value)
	}

	res.Body = fdp.ParseBody(req, res, res.Body, res.BodyAppend)
}

//ParseBody parses body respecting bodyAppend and replacing variables from request
func (fdp FakeDataParse) ParseBody(req *definition.Request, res *definition.Response, body string, bodyAppend string) string {
	resultBody := fdp.ReplaceVars(req, res, body)
	if bodyAppend != "" {
		resultBodyAppend := fdp.ReplaceVars(req, res, bodyAppend)

		if isJSON(resultBody) && isJSON(resultBodyAppend) {
			resultBody = joinJSON(resultBody, resultBodyAppend)
		} else if isJSON(resultBody) && !isJSON(resultBodyAppend) {
			// strip resultBodyAppend as it is not in appropriate format
			log.Printf("BodyAppend not in JSON format : %s\n", resultBodyAppend)
		} else {
			resultBody += resultBodyAppend
		}
	}

	return resultBody
}

func joinJSON(inputs ...string) string {
	if len(inputs) == 1 {
		return inputs[0]
	}

	result := gabs.New()
	for _, input := range inputs {
		jsonParsed, _ := gabs.ParseJSON([]byte(input))
		children, _ := jsonParsed.S().ChildrenMap()

		for key, child := range children {
			result.Set(child.Data(), key)
		}
	}

	return result.String()
}

func getStringPart(input string, pattern string, groupName string) (string, bool) {
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

func getQueryStringParam(req *definition.Request, name string) (string, bool) {

	if len(req.QueryStringParameters) == 0 {
		return "", false
	}
	value, f := req.QueryStringParameters[name]
	if !f {
		return "", false
	}

	return value[0], true
}

func getCookieParam(req *definition.Request, name string) (string, bool) {

	if len(req.Cookies) == 0 {
		return "", false
	}
	value, f := req.Cookies[name]
	if !f {
		return "", false
	}

	return value, true
}
