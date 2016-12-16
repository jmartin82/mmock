package vars

import (
	"testing"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/vars/fakedata"
)

func getProcessor() VarsProcessor {
	return VarsProcessor{FillerFactory: MockFillerFactory{}, FakeAdapter: fakedata.NewDummyDataFaker("AleixMG")}
}

func TestReplaceTags(t *testing.T) {

	req := definition.Request{}
	val := make(definition.Values)
	val["param1"] = []string{"valParam"}
	req.QueryStringParameters = val

	cookie := make(definition.Cookies)
	cookie["cookie1"] = "valCookie"
	req.Cookies = cookie

	res := definition.Response{}
	res.Body = "Request {{request.query.param1}}. Cookie: {{request.cookie.cookie1}}. Random: {{fake.UserName}}"

	cookie = make(definition.Cookies)
	cookie["cookie1"] = "valCookie"
	cookie["cookie2"] = "{{fake.UserName}}"
	res.Cookies = cookie

	val = make(definition.Values)
	val["header1"] = []string{"valHeader"}
	val["header2"] = []string{"valHeader", "{{request.query.param1}}"}

	res.Headers = val

	mock := definition.Mock{Request: req, Response: res}
	varsProcessor := getProcessor("testData")
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != "Request valParam. Cookie: valCookie. Random: AleixMG" {
		t.Error("Replaced tags in body not match", res.Body)
	}

	if mock.Response.Cookies["cookie2"] != "AleixMG" {
		t.Error("Replaced tags in cookie match", mock.Response.Cookies["cookie2"])
	}

	if mock.Response.Headers["header2"][1] != "valParam" {
		t.Error("Replaced tags in headers match", mock.Response.Headers["header2"][1])
	}
}

func TestReplaceUndefinedFakeTag(t *testing.T) {
	req := definition.Request{}
	val := make(definition.Values)
	val["param1"] = []string{"valParam"}
	req.QueryStringParameters = val

	cookie := make(definition.Cookies)
	cookie["cookie1"] = "valCookie"
	req.Cookies = cookie

	res := definition.Response{}
	res.Body = "Request {{request.query.param2}}. Cookie: {{request.cookie.cookie2}}. Random: {{fake.otherOption}}"

	mock := definition.Mock{Request: req, Response: res}
	varsProcessor := getProcessor("testData")
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != "Request {{request.query.param2}}. Cookie: {{request.cookie.cookie2}}. Random: {{fake.otherOption}}" {
		t.Error("Replaced tags in body not match", mock.Response.Body)
	}

}

func TestReplaceTagWithSpace(t *testing.T) {
	req := definition.Request{}
	val := make(definition.Values)
	val["param1"] = []string{"valParam"}
	req.QueryStringParameters = val

	cookie := make(definition.Cookies)
	cookie["cookie1"] = "valCookie"
	req.Cookies = cookie

	res := definition.Response{}
	res.Body = "Request {{ request.query.param1}}. Cookie: {{request.cookie.cookie1 }}. Random: {{fake.UserName }}"

	mock := definition.Mock{Request: req, Response: res}
	varsProcessor := getProcessor("testData")
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != "Request valParam. Cookie: valCookie. Random: AleixMG" {
		t.Error("Replaced tags in body not match", mock.Response.Body)
	}
}

func TestReplaceUrlRegex(t *testing.T) {
	req := definition.Request{}
	res := definition.Response{}

	req.Path = "/users/15"
	res.Body = "{ \"id\": {{request.url./users/(?P<value>\\d+)}} }"

	mock := definition.Mock{Request: req, Response: res}
	varsProcessor := getProcessor("testData")
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != "{ \"id\": 15 }" {
		t.Error("Replaced url regex in body not match", mock.Response.Body)
	}
}

func TestReplaceBodyRegex(t *testing.T) {
	req := definition.Request{}
	res := definition.Response{}

	req.Path = "/"
	req.Body = "/users/15"
	res.Body = "{ \"id\": {{request.body.users/(?P<value>\\d+)}} }"

	mock := definition.Mock{Request: req, Response: res}
	varsProcessor := getProcessor("testData")
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != "{ \"id\": 15 }" {
		t.Error("Replaced body regex in body not match", mock.Response.Body)
	}
}

func TestReplaceTagWithParameter(t *testing.T) {
	req := definition.Request{}

	res := definition.Response{}
	res.Body = "Random: {{fake.CharactersN(15)}}"

	mock := definition.Mock{Request: req, Response: res}
	varsProcessor := getProcessor("testData")
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != "Random: AleixMG15" {
		t.Error("Replaced tags in body not match", mock.Response.Body)
	}
}

func TestReplaceTagWithParameterNoParameterPassed(t *testing.T) {
	req := definition.Request{}

	res := definition.Response{}
	res.Body = "Random: {{fake.CharactersN}}"

	mock := definition.Mock{Request: req, Response: res}
	varsProcessor := getProcessor("testData")
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != "Random: {{fake.CharactersN}}" {
		t.Error("Replaced tags in body not match", mock.Response.Body)
	}
}
