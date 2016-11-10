package parse

import (
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/parse/fakedata"
)

//FakeDataParse parses the data looking for fake data tags or request data tags
type FakeDataParse struct {
	Fake fakedata.DataFaker
}

func (fdp FakeDataParse) getQueryStringParam(req *definition.Request, name string) (string, bool) {

	if len(req.QueryStringParameters) == 0 {
		return "", false
	}
	value, f := req.QueryStringParameters[name]
	if !f {
		return "", false
	}

	return value[0], true
}

func (fdp FakeDataParse) getCookieParam(req *definition.Request, name string) (string, bool) {

	if len(req.Cookies) == 0 {
		return "", false
	}
	value, f := req.Cookies[name]
	if !f {
		return "", false
	}

	return value, true
}

func (fdp FakeDataParse) getURLPart(req *definition.Request, pattern string) (string, bool) {
	r, error := regexp.Compile(pattern)
	if error != nil{
		return "", false
	} 

	match := r.FindStringSubmatch(req.Path)
	result := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			result[name] = match[i]
		}
	}

	value, present := result["Value"]

	return value, present
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

func (fdp FakeDataParse) replaceVars(req *definition.Request, input string) string {
	r := regexp.MustCompile(`\{\{\s*([^\}]+)\s*\}\}`)

	return r.ReplaceAllStringFunc(input, func(raw string) string {
		found := false
		s := ""
		tag := strings.Trim(raw[2:len(raw)-2], " ")
		if tag == "request.body" {
			s = req.Body
			found = true
		} else if i := strings.Index(tag, "request.url."); i == 0 {
			s, found = fdp.getURLPart(req, tag[12:])
		} else if i := strings.Index(tag, "request.query."); i == 0 {
			s, found = fdp.getQueryStringParam(req, tag[14:])
		} else if i := strings.Index(tag, "request.cookie."); i == 0 {
			s, found = fdp.getCookieParam(req, tag[15:])
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
			res.Headers[header][i] = fdp.replaceVars(req, value)
		}

	}
	for cookie, value := range res.Cookies {
		res.Cookies[cookie] = fdp.replaceVars(req, value)
	}

	res.Body = fdp.replaceVars(req, res.Body)
}
