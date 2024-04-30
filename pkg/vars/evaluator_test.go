package vars

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/jmartin82/mmock/v3/pkg/mock"

	"strconv"

	"strings"
)

// DummyDataFaker is used in tests
type DummyDataFaker struct {
	Dummy string
}

func NewDummyDataFaker(dummyString string) DummyDataFaker {
	result := DummyDataFaker{Dummy: dummyString}
	return result
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
func (ddf DummyDataFaker) CharactersN(n int) string {
	return ddf.Dummy + strconv.Itoa(n)
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
func (ddf DummyDataFaker) DigitsN(n int) string {
	return ddf.Dummy + strconv.Itoa(n)
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
func (ddf DummyDataFaker) MonthNum() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Paragraph() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Paragraphs() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) ParagraphsN(n int) string {
	return ddf.Dummy + strconv.Itoa(n)
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
func (ddf DummyDataFaker) SentencesN(n int) string {
	return ddf.Dummy + strconv.Itoa(n)
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
func (ddf DummyDataFaker) WordsN(n int) string {
	return ddf.Dummy + strconv.Itoa(n)
}
func (ddf DummyDataFaker) Zip() string {
	return ddf.Dummy
}
func (ddf DummyDataFaker) Int(n int) string {
	return ddf.Dummy + strconv.Itoa(n)
}
func (ddf DummyDataFaker) IntMinMax(values ...int) string {
	return ddf.Dummy + strconv.Itoa(values[0]+1)
}
func (ddf DummyDataFaker) Float(n int) string {
	return ddf.Dummy + strconv.Itoa(n)
}
func (ddf DummyDataFaker) UUID() string {
	return "00000000-0000-0000-0000-000000000000"
}
func (ddf DummyDataFaker) Hex(n int) string {
	return strings.Repeat("0", n)
}

func getProcessor() ResponseMessageEvaluator {
	dfp := NewDummyDataFaker("AleixMG")
	ff := NewFillerFactory(dfp)
	return ResponseMessageEvaluator{FillerFactory: ff}
}

func TestReplaceTags(t *testing.T) {

	req := mock.Request{}
	req.Body = "hi!"

	val := make(mock.Values)
	val["param1"] = []string{"valParam"}

	req.QueryStringParameters = val

	cookie := make(mock.Cookies)
	cookie["cookie1"] = "valCookie"
	req.Cookies = cookie

	res := mock.Response{}
	cb := mock.Callback{}
	res.Body = "Request Body {{request.body}}. Query {{request.query.param1}}. Cookie: {{request.cookie.cookie1}}. Random: {{fake.UserName}}"
	cb.Body = "Callback Body {{request.body}}. Query {{request.query.param1}}. Cookie: {{request.cookie.cookie1}}. Random: {{fake.UserName}}"

	cookie = make(mock.Cookies)
	cookie["cookie1"] = "valCookie"
	cookie["cookie2"] = "{{fake.UserName}}"

	res.Cookies = cookie
	cb.Cookies = cookie

	val = make(mock.Values)
	val["header1"] = []string{"valHeader"}
	val["header2"] = []string{"valHeader", "{{request.query.param1}}"}

	res.Headers = val
	cb.Headers = val

	mock := mock.Definition{Request: req, Response: res, Callback: cb}
	varsProcessor := getProcessor()
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != "Request Body hi!. Query valParam. Cookie: valCookie. Random: AleixMG" {
		t.Error("Replaced tags in body not match", res.Body)
	}

	if mock.Response.Cookies["cookie2"] != "AleixMG" {
		t.Error("Replaced tags in cookie match", mock.Response.Cookies["cookie2"])
	}

	if mock.Response.Headers["header2"][1] != "valParam" {
		t.Error("Replaced tags in headers match", mock.Response.Headers["header2"][1])
	}

	if mock.Callback.Cookies["cookie2"] != "AleixMG" {
		t.Error("Replaced tags in Callback cookie match", mock.Callback.Cookies["cookie2"])
	}

	if mock.Callback.Headers["header2"][1] != "valParam" {
		t.Error("Replaced tags in Callback headers match", mock.Callback.Headers["header2"][1])
	}

	if mock.Callback.Body != "Callback Body hi!. Query valParam. Cookie: valCookie. Random: AleixMG" {
		t.Error("Replaced tags in body not match", cb.Body)
	}
}

func TestReplaceUndefinedFakeTag(t *testing.T) {
	req := mock.Request{}
	val := make(mock.Values)
	val["param1"] = []string{"valParam"}
	req.QueryStringParameters = val

	cookie := make(mock.Cookies)
	cookie["cookie1"] = "valCookie"
	req.Cookies = cookie

	res := mock.Response{}
	res.Body = "Request {{request.query.param2}}. Cookie: {{request.cookie.cookie2}}. Random: {{fake.otherOption}}"

	mock := mock.Definition{Request: req, Response: res}
	varsProcessor := getProcessor()
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != "Request {{request.query.param2}}. Cookie: {{request.cookie.cookie2}}. Random: {{fake.otherOption}}" {
		t.Error("Replaced tags in body not match", mock.Response.Body)
	}

}

func TestReplaceTagWithSpace(t *testing.T) {
	req := mock.Request{}
	val := make(mock.Values)
	val["param1"] = []string{"valParam"}
	req.QueryStringParameters = val

	cookie := make(mock.Cookies)
	cookie["cookie1"] = "valCookie"
	req.Cookies = cookie

	res := mock.Response{}
	res.Body = "Request {{ request.query.param1}}. Cookie: {{request.cookie.cookie1 }}. Random: {{fake.UserName }}"

	mock := mock.Definition{Request: req, Response: res}
	varsProcessor := getProcessor()
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != "Request valParam. Cookie: valCookie. Random: AleixMG" {
		t.Error("Replaced tags in body not match", mock.Response.Body)
	}
}

func TestReplaceUrlPathVars(t *testing.T) {

	mockReq := mock.Request{}
	mockReq.Path = "/users/:id"
	res := mock.Response{}
	res.Body = "{ \"id\": {{request.path.id}} }"

	m := mock.Definition{Request: mockReq, Response: res}
	varsProcessor := getProcessor()

	req := mock.Request{}
	req.Path = "/users/15"
	varsProcessor.Eval(&req, &m)

	if m.Response.Body != "{ \"id\": 15 }" {
		t.Error("Replaced url param in body not match", m.Response.Body)
	}
}

func TestReplaceTagWithParameter(t *testing.T) {
	req := mock.Request{}

	res := mock.Response{}
	res.Body = "Random: {{fake.CharactersN(15)}}"

	m := mock.Definition{Request: req, Response: res}
	varsProcessor := getProcessor()
	varsProcessor.Eval(&req, &m)

	if m.Response.Body != "Random: AleixMG15" {
		t.Error("Replaced tags in body not match", m.Response.Body)
	}
}

func TestReplaceTagWithParameterNoParameterPassed(t *testing.T) {
	req := mock.Request{}

	res := mock.Response{}
	res.Body = "Random: {{fake.CharactersN}}"

	mock := mock.Definition{Request: req, Response: res}
	varsProcessor := getProcessor()
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != "Random: {{fake.CharactersN}}" {
		t.Error("Replaced tags in body not match", mock.Response.Body)
	}
}

func TestReplaceMissingTags(t *testing.T) {
	req := mock.Request{}

	res := mock.Response{}
	res.Body = "Request Body {{request.body}}. Query {{request.query.param1}}. Cookie: {{request.cookie.cookie1}}."

	mock := mock.Definition{Request: req, Response: res}
	varsProcessor := getProcessor()
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != "Request Body {{request.body}}. Query {{request.query.param1}}. Cookie: {{request.cookie.cookie1}}." {
		t.Error("Replaced missing tags not match", mock.Response.Body)
	}
}

func TestReplaceFormUrlEncodedBodyTags(t *testing.T) {
	req := mock.Request{}
	req.Body = "one=foo&two[array]=bar"
	req.Headers = make(mock.Values)
	req.Headers["Content-Type"] = []string{"application/x-www-form-urlencoded"}

	res := mock.Response{}
	res.Body = "Form data placeholders. One '{{request.body.one}}'. Two '{{request.body.two[array]}}'."

	mock := mock.Definition{Request: req, Response: res}
	varsProcessor := getProcessor()
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != "Form data placeholders. One 'foo'. Two 'bar'." {
		t.Error("Replaced tags from body form do not match", mock.Response.Body)
	}
}

func TestReplaceUrlInfo(t *testing.T) {

	req := mock.Request{}
	req.Path = "/home"

	val := make(mock.Values)
	val["param1"] = []string{"valParam1", "valParam2"}
	val["param2"] = []string{"valParam1"}
	req.QueryStringParameters = val

	req.Scheme = "ws"
	req.Host = "example.com"
	req.Port = "8001"
	req.Fragment = "anchor"

	res := mock.Response{}
	res.Body = "{{request.scheme}}://{{request.hostname}}:{{request.port}}{{request.path}}#{{request.fragment}}"

	m := mock.Definition{Request: req, Response: res}
	varsProcessor := getProcessor()
	varsProcessor.Eval(&req, &m)

	if m.Response.Body != "ws://example.com:8001/home#anchor" {
		t.Error("Replaced url info from body do not match", m.Response.Body)
	}

	res = mock.Response{}
	res.Body = "{{request.url}}"

	m = mock.Definition{Request: req, Response: res}
	varsProcessor.Eval(&req, &m)

	if m.Response.Body != "ws://example.com:8001/home?param1=valParam1&param1=valParam2&param2=valParam1#anchor" {
		t.Error("Replaced url info from body do not match", m.Response.Body)
	}
}

func TestReplaceJsonBodyEncodedTags(t *testing.T) {
	req := mock.Request{}
	req.Headers = make(mock.Values)
	req.Headers["Content-Type"] = []string{"application/json"}
	req.Body = `
{
  "email": "hilari@hilarimoragrega.com",
  "age": 34,
  "height": 5.66,
  "weight": null,
  "level": -8,
  "active": true,
  "friends": [
    "jordi.martin@gmail.com", 
    "alfons.faubert@gmail.com"
  ],
  "attributes": {
    "programming": 15,
    "trolling": 27
  },
  "tracking": {
    "uuid":"0bd74115-2307-458f-8288-b726724045ef",
    "nesting": {
      "level": "nesting is ok"
    },
    "discarded": "do not return"
  }
}
`
	res := mock.Response{}
	res.Body = `
{
  "email": "{{request.body.email}}",
  "age": {{request.body.age}},
  "height": {{request.body.height}},
  "weight": {{request.body.weight}},
  "active": {{request.body.active}},
  "level": {{request.body.level}},
  "friends": {{request.body.friends}},
  "first-friend": "{{request.body.friends.0}}",
  "last-friend": "{{request.body.friends.1}}",
  "attributes": {{request.body.attributes}},
  "tracking": {
    "uuid": "{{request.body.tracking.uuid}}",
    "deeper": {
      "level": "{{request.body.tracking.nesting.level}}"
    }
  }
}
`

	expected := `
{
  "email": "hilari@hilarimoragrega.com",
  "age": 34,
  "height": 5.66,
  "weight": ,
  "active": true,
  "level": -8,
  "friends": [
    "jordi.martin@gmail.com", 
    "alfons.faubert@gmail.com"
  ],
  "first-friend": "jordi.martin@gmail.com",
  "last-friend": "alfons.faubert@gmail.com",
  "attributes": {
    "programming": 15,
    "trolling": 27
  },
  "tracking": {
    "uuid": "0bd74115-2307-458f-8288-b726724045ef",
    "deeper": {
      "level": "nesting is ok"
    }
  }
}
`
	mock := mock.Definition{Request: req, Response: res}
	varsProcessor := getProcessor()
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != expected {
		t.Error("Replaced tags from body form do not match", mock.Response.Body)
	}
}

func TestReplaceXmlBodyEncodedTags(t *testing.T) {
	req := mock.Request{}
	req.Headers = make(mock.Values)
	req.Headers["Content-Type"] = []string{"application/xml"}
	req.Body = `
<?xml version="1.0" encoding="UTF-8"?>
<root>
	<active>true</active>
	<age>34</age>
	<attributes>
		<programming>15</programming>
		<trolling>27</trolling>
	</attributes>
	<email>hilari@hilarimoragrega.com</email>
	<friends>
		<element>jordi.martin@gmail.com</element>
		<element>alfons.faubert@gmail.com</element>
	</friends>
	<height>5.66</height>
	<level>-8</level>
	<tracking>
		<discarded>do not return</discarded>
		<nesting>
			<level>nesting is ok</level>
		</nesting>
		<uuid>0bd74115-2307-458f-8288-b726724045ef</uuid>
	</tracking>
</root>
`
	res := mock.Response{}
	res.Body = `
<?xml version="1.0" encoding="UTF-8"?>
<root>
	<active>{{request.root.body.active}}</active>
	<age>{{request.body.root.email}}</age>
	<email>{{request.body.root.email}}</email>
	<friends>
		<element>{{request.body.root.friends.element.0}}</element>
		<element>{{request.body.root.friends.element.1}}</element>
	</friends>
	<level>{{request.body.root.level}}</level>
	<tracking>
		<nesting>
			<level>{{request.body.root.tracking.nesting.level}}</level>
		</nesting>
		<uuid>{{request.body.root.tracking.uuid}}</uuid>
	</tracking>
</root>
`

	expected := `
<?xml version="1.0" encoding="UTF-8"?>
<root>
	<active>{{request.root.body.active}}</active>
	<age>hilari@hilarimoragrega.com</age>
	<email>hilari@hilarimoragrega.com</email>
	<friends>
		<element>jordi.martin@gmail.com</element>
		<element>alfons.faubert@gmail.com</element>
	</friends>
	<level>-8</level>
	<tracking>
		<nesting>
			<level>nesting is ok</level>
		</nesting>
		<uuid>0bd74115-2307-458f-8288-b726724045ef</uuid>
	</tracking>
</root>
`
	mock := mock.Definition{Request: req, Response: res}
	varsProcessor := getProcessor()
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != expected {
		t.Error("Replaced tags from body form do not match", mock.Response.Body)
	}
}

func TestReplaceTagsCallback(t *testing.T) {

	req := mock.Request{}
	req.Body = "hi Isona!"
	val := make(mock.Values)
	val["param1"] = []string{"valParam"}
	req.QueryStringParameters = val

	cookie := make(mock.Cookies)
	cookie["cookie1"] = "valCookie"
	req.Cookies = cookie

	cb := mock.Callback{}
	cb.Body = "Request Body {{request.body}}. Query {{request.query.param1}}. Cookie: {{request.cookie.cookie1}}. Random: {{fake.UserName}}"

	mock := mock.Definition{Request: req, Callback: cb}
	varsProcessor := getProcessor()
	varsProcessor.Eval(&req, &mock)

	if mock.Callback.Body != "Request Body hi Isona!. Query valParam. Cookie: valCookie. Random: AleixMG" {
		t.Error("Replaced tags in callback body not match", cb.Body)
	}
}

func TestReplaceBigFile(t *testing.T) {
	content := []byte("{{request.body}} this is a big file with holders replaced")
	dir, err := ioutil.TempDir("", "mmock")
	if err != nil {
		t.Errorf("Error creating temporary folder")
	}

	tmpfn := filepath.Join(dir, "bigfile")
	if err := ioutil.WriteFile(tmpfn, content, 0666); err != nil {
		t.Errorf("Error updating temporary file")
	}

	defer os.RemoveAll(dir) // clean up

	req := mock.Request{}
	req.Body = "hi! Isona."

	res := mock.Response{}
	res.Body = fmt.Sprintf("Big file: {{file.contents(%s)}}", tmpfn)

	mock := mock.Definition{Request: req, Response: res}
	varsProcessor := getProcessor()
	varsProcessor.Eval(&req, &mock)

	if mock.Response.Body != "Big file: hi! Isona. this is a big file with holders replaced" {
		t.Error("Replaced tags in a external stream doesn't work.", res.Body, mock.Response.Body)
	}
}
