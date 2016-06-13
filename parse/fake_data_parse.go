package parse

import (
	"log"
	"mmock/definition"
	"mmock/parse/fakedata"
	"reflect"
	"regexp"
	"strings"
)

type FakeDataParse struct {
	Fake fakedata.DataFaker
}

func (this FakeDataParse) getQueryStringParam(req *definition.Request, name string) (string, bool) {

	if len(req.QueryStringParameters) == 0 {
		return "", false
	}
	value, f := req.QueryStringParameters[name]
	if !f {
		return "", false
	}

	return value[0], true
}

func (this FakeDataParse) getCookieParam(req *definition.Request, name string) (string, bool) {

	if len(req.Cookies) == 0 {
		return "", false
	}
	value, f := req.Cookies[name]
	if !f {
		return "", false
	}

	return value, true
}

func (this FakeDataParse) call(data reflect.Value, name string) string {
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

func (this FakeDataParse) callMethod(name string) (string, bool) {
	found := false
	data := reflect.ValueOf(this.Fake)
	typ := data.Type()
	if nMethod := data.Type().NumMethod(); nMethod > 0 {
		for i := 0; i < nMethod; i++ {
			method := typ.Method(i)
			if strings.ToLower(method.Name) == strings.ToLower(name) {
				found = true // we found the name regardless
				// does receiver type match? (pointerness might be off)
				if typ == method.Type.In(0) {
					return this.call(data, method.Name), found
				}
			}
		}
	}
	return "", found
}

func (this FakeDataParse) replaceVars(req *definition.Request, input string) string {
	r := regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_\.]+)\s*\}\}`)

	return r.ReplaceAllStringFunc(input, func(raw string) string {
		found := false
		s := ""
		tag := strings.Trim(raw[2:len(raw)-2], " ")
		if i := strings.Index(tag, "request.query."); i == 0 {
			s, found = this.getQueryStringParam(req, tag[14:len(tag)])
		} else if i := strings.Index(tag, "request.cookie."); i == 0 {
			s, found = this.getCookieParam(req, tag[15:len(tag)])
		} else if i := strings.Index(tag, "fake."); i == 0 {
			s, found = this.callMethod(tag[5:len(tag)])
		}

		if !found {
			log.Printf("Defined tag {{%s}} not found\n", tag)
			return raw
		}
		return s
	})
}

func (this FakeDataParse) Parse(req *definition.Request, res *definition.Response) {
	for header, values := range res.Headers {
		for i, value := range values {
			res.Headers[header][i] = this.replaceVars(req, value)
		}

	}
	for cookie, value := range res.Cookies {
		res.Cookies[cookie] = this.replaceVars(req, value)
	}

	res.Body = this.replaceVars(req, res.Body)
}
