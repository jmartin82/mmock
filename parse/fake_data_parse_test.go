package parse

import (
	"testing"

	"github.com/jmartin82/mmock/definition"
)

type DummyDataFaker struct {
	Dummy string
}

func (ddf DummyDataFaker) Brand() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Character() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Characters() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) City() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Color() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Company() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Continent() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Country() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) CreditCardVisa() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) CreditCardMasterCard() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) CreditCardAmericanExpress() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Currency() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) CurrencyCode() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Day() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Digits() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) EmailAddress() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) FirstName() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) FullName() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) LastName() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Gender() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) IPv4() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Language() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Model() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Month() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Year() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) MonthShort() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Paragraph() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Paragraphs() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Phone() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Product() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Sentence() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Sentences() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) SimplePassword() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) State() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) StateAbbrev() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Street() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) StreetAddress() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) UserName() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) WeekDay() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Word() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Words() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Zip() string {
	return ddf.Dummy
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
