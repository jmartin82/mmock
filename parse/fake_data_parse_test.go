package parse

import (
	"mmock/definition"
	"testing"
)

type DummyDataFaker struct {
	Dummy string
}

func (this DummyDataFaker) Brand() string {
	return this.Dummy
}
func (this DummyDataFaker) Character() string {
	return this.Dummy
}
func (this DummyDataFaker) Characters() string {
	return this.Dummy
}
func (this DummyDataFaker) City() string {
	return this.Dummy
}
func (this DummyDataFaker) Color() string {
	return this.Dummy
}
func (this DummyDataFaker) Company() string {
	return this.Dummy
}
func (this DummyDataFaker) Continent() string {
	return this.Dummy
}
func (this DummyDataFaker) Country() string {
	return this.Dummy
}
func (this DummyDataFaker) CreditCardVisa() string {
	return this.Dummy
}
func (this DummyDataFaker) CreditCardMasterCard() string {
	return this.Dummy
}
func (this DummyDataFaker) CreditCardAmericanExpress() string {
	return this.Dummy
}
func (this DummyDataFaker) Currency() string {
	return this.Dummy
}
func (this DummyDataFaker) CurrencyCode() string {
	return this.Dummy
}
func (this DummyDataFaker) Day() string {
	return this.Dummy
}
func (this DummyDataFaker) Digits() string {
	return this.Dummy
}
func (this DummyDataFaker) EmailAddress() string {
	return this.Dummy
}
func (this DummyDataFaker) FirstName() string {
	return this.Dummy
}
func (this DummyDataFaker) FullName() string {
	return this.Dummy
}
func (this DummyDataFaker) LastName() string {
	return this.Dummy
}
func (this DummyDataFaker) Gender() string {
	return this.Dummy
}
func (this DummyDataFaker) IPv4() string {
	return this.Dummy
}
func (this DummyDataFaker) Language() string {
	return this.Dummy
}
func (this DummyDataFaker) Model() string {
	return this.Dummy
}
func (this DummyDataFaker) Month() string {
	return this.Dummy
}
func (this DummyDataFaker) Year() string {
	return this.Dummy
}
func (this DummyDataFaker) MonthShort() string {
	return this.Dummy
}
func (this DummyDataFaker) Paragraph() string {
	return this.Dummy
}
func (this DummyDataFaker) Paragraphs() string {
	return this.Dummy
}
func (this DummyDataFaker) Phone() string {
	return this.Dummy
}
func (this DummyDataFaker) Product() string {
	return this.Dummy
}
func (this DummyDataFaker) Sentence() string {
	return this.Dummy
}
func (this DummyDataFaker) Sentences() string {
	return this.Dummy
}
func (this DummyDataFaker) SimplePassword() string {
	return this.Dummy
}
func (this DummyDataFaker) State() string {
	return this.Dummy
}
func (this DummyDataFaker) StateAbbrev() string {
	return this.Dummy
}
func (this DummyDataFaker) Street() string {
	return this.Dummy
}
func (this DummyDataFaker) StreetAddress() string {
	return this.Dummy
}
func (this DummyDataFaker) UserName() string {
	return this.Dummy
}
func (this DummyDataFaker) WeekDay() string {
	return this.Dummy
}
func (this DummyDataFaker) Word() string {
	return this.Dummy
}
func (this DummyDataFaker) Words() string {
	return this.Dummy
}
func (this DummyDataFaker) Zip() string {
	return this.Dummy
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
