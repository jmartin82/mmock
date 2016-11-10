package parse

import (
	"log"
	"reflect"
	"regexp"
	"strings"
	"encoding/json"

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
			s, found = req.GetURLPart(tag[12:], "Value")
		} else if i := strings.Index(tag, "request.query."); i == 0 {
			s, found = req.GetQueryStringParam(tag[14:])
		} else if i := strings.Index(tag, "request.cookie."); i == 0 {
			s, found = req.GetCookieParam(tag[15:])
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

	res.Body = fdp.ParseBody(res.Body, res.BodyAppend, req)

}

//ParseBody parses body respecting bodyAppend and replacing variables from request
func (fdp FakeDataParse) ParseBody(body string, bodyAppend string, req *definition.Request) string{
	resultBody := fdp.replaceVars(req, body)
	if(bodyAppend != ""){
		resultBodyAppend := fdp.replaceVars(req, bodyAppend)

		if isJSON(resultBody) && isJSON(resultBodyAppend){
			resultBody = fdp.JoinJSON(resultBody, resultBodyAppend)
		} else if isJSON(resultBody) && !isJSON(resultBodyAppend){
			// strip resultBodyAppend as it is not in appropriate format
			log.Printf("BodyAppend not in JSON format : %s\n", resultBodyAppend)
		} else {
			resultBody += resultBodyAppend
		}
	}

	return resultBody
}

//JoinJSON joins the properties of the passed jsons 
func (fdp FakeDataParse) JoinJSON(inputs ...string) string {
	if len(inputs) == 1 {
		return inputs[0]
	} 

	result := gabs.New()
	for _, input := range inputs{
		jsonParsed, _ := gabs.ParseJSON([]byte(input))
		children, _ := jsonParsed.S().ChildrenMap()

		for key, child := range children {
			result.Set(child.Data(), key)
		}
	}

	return result.String()
}
