package vars

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/vars/fakedata"
)

//FakeVars parses the data looking for fake data tags or request data tags
type FakeVars struct {
	Fake fakedata.DataFaker
}

func (fdp FakeVars) call(data reflect.Value, name string) string {
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

func (fdp FakeVars) callMethod(name string) (string, bool) {
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

func (fdp FakeVars) Fill(m *definition.Mock, input string) string {
	r := regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_\.]+)\s*\}\}`)

	return r.ReplaceAllStringFunc(input, func(raw string) string {
		found := false
		s := ""
		tag := strings.Trim(raw[2:len(raw)-2], " ")
		if i := strings.Index(tag, "fake."); i == 0 {
			s, found = fdp.callMethod(tag[5:])
		}

		if !found {
			return raw
		}
		return s
	})
}
