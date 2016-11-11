package parse

import (
	"testing"

	"github.com/jmartin82/mmock/definition"
)

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

	faker := FakeDataParse{DummyDataFaker{"AleixMG"}}
	faker.Parse(&req, &res)

	if res.Body != "Request valParam. Cookie: valCookie. Random: AleixMG" {
		t.Error("Replaced tags in body not match", res.Body)
	}

	if res.Cookies["cookie2"] != "AleixMG" {
		t.Error("Replaced tags in cookie match", res.Cookies["cookie2"])
	}

	if res.Headers["header2"][1] != "valParam" {
		t.Error("Replaced tags in headers match", res.Headers["header2"][1])
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

	faker := FakeDataParse{DummyDataFaker{"AleixMG"}}
	faker.Parse(&req, &res)

	if res.Body != "Request {{request.query.param2}}. Cookie: {{request.cookie.cookie2}}. Random: {{fake.otherOption}}" {
		t.Error("Replaced tags in body not match", res.Body)
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
	res.Body = "Request {{request.query.param1}}. Cookie: {{request.cookie.cookie1}}. Random: {{fake.UserName}}"

	faker := FakeDataParse{DummyDataFaker{"AleixMG"}}
	faker.Parse(&req, &res)

	if res.Body != "Request valParam. Cookie: valCookie. Random: AleixMG" {
		t.Error("Replaced tags in body not match", res.Body)
	}
}

func TestReplaceUrlRegex(t *testing.T) {
	req := definition.Request{}
	res := definition.Response{}

	req.Path = "/users/15"
	res.Body = "{ \"id\": {{request.url./users/(?P<value>\\d+)}} }"

	faker := FakeDataParse{DummyDataFaker{"AleixMG"}}
	faker.Parse(&req, &res)

	if res.Body != "{ \"id\": 15 }" {
		t.Error("Replaced url regex in body not match", res.Body)
	}
}

func TestReplaceBodyRegex(t *testing.T) {
	req := definition.Request{}
	res := definition.Response{}

	req.Path = "/"
	req.Body = "/users/15"
	res.Body = "{ \"id\": {{request.body.users/(?P<value>\\d+)}} }"

	faker := FakeDataParse{DummyDataFaker{"AleixMG"}}
	faker.Parse(&req, &res)

	if res.Body != "{ \"id\": 15 }" {
		t.Error("Replaced body regex in body not match", res.Body)
	}
}

func TestBodyAppendNoJson(t *testing.T) {
	req := definition.Request{}
	res := definition.Response{}

	req.Path = "/users/15"
	res.Body = "Test"
	res.BodyAppend = "Append"

	faker := FakeDataParse{DummyDataFaker{"AleixMG"}}
	faker.Parse(&req, &res)

	if res.Body != "TestAppend" {
		t.Error("BodyAppend text not added as expected", res.Body)
	}
}

func TestBodyAppendBodyJsonAppendNoJson(t *testing.T) {
	req := definition.Request{}
	res := definition.Response{}

	req.Path = "/users/15"
	res.Body = "{\"body\":1}"
	res.BodyAppend = "Append"

	faker := FakeDataParse{DummyDataFaker{"AleixMG"}}
	faker.Parse(&req, &res)

	if res.Body != "{\"body\":1}" {
		t.Error("BodyAppend should be discarded as the the body is in JSON format, but the BodyAppend not", res.Body)
	}
}

func TestBodyAppendBodyJsonAppend(t *testing.T) {
	req := definition.Request{}
	res := definition.Response{}

	req.Path = "/users/15"
	res.Body = "{\"body\":1}"
	res.BodyAppend = "{\"append\":2}"

	faker := FakeDataParse{DummyDataFaker{"AleixMG"}}
	faker.Parse(&req, &res)

	if res.Body != "{\"append\":2,\"body\":1}" {
		t.Error("BodyAppend fields shoud be added to the response body", res.Body)
	}
}

func TestBodyAppendBodyJsonAppendOverrideValue(t *testing.T) {
	req := definition.Request{}
	res := definition.Response{}

	req.Path = "/users/15"
	res.Body = "{\"body\":1}"
	res.BodyAppend = "{\"body\":2}"

	faker := FakeDataParse{DummyDataFaker{"AleixMG"}}
	faker.Parse(&req, &res)

	if res.Body != "{\"body\":2}" {
		t.Error("BodyAppend fields should override the response body fields", res.Body)
	}
}
