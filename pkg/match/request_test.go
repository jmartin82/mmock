package match

import (
	"github.com/jmartin82/mmock/v3/pkg/match/payload"
	"testing"

	"github.com/jmartin82/mmock/v3/pkg/mock"
)

func TestMatchMethod(t *testing.T) {
	req := mock.Request{}
	req.Method = "GET"

	m := mock.Definition{}
	m.Request.Method = "GET"

	mm := Request{}
	if b, err := mm.Match(&req, &m, true); !b {
		t.Error(err)
	}

	req.Method = "POST"
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchScheme(t *testing.T) {
	req := mock.Request{}
	req.Scheme = "https"

	m := mock.Definition{}
	m.Request.Scheme = "https"

	mm := Request{}
	if b, err := mm.Match(&req, &m, true); !b {
		t.Error(err)
	}

	req.Scheme = "http"
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchFragment(t *testing.T) {
	req := mock.Request{}
	req.Fragment = "fragment"

	m := mock.Definition{}
	m.Request.Fragment = "fragment"

	mm := Request{}
	if b, err := mm.Match(&req, &m, true); !b {
		t.Error(err)
	}

	req.Fragment = "nothing"
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchPath(t *testing.T) {

	req := mock.Request{}
	req.Path = "/a/b/c"

	m := mock.Definition{}
	m.Request.Path = "/a/b/c"

	mm := Request{}
	if b, err := mm.Match(&req, &m, true); !b {
		t.Error(err)
	}

	req.Path = "/a/b/d"
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}
}

func TestPathVars(t *testing.T) {
	req := mock.Request{}
	req.Path = "/a/b/c"

	m := mock.Definition{}
	m.Request.Path = "/a/:b/:c"

	mm := Request{}

	if b, err := mm.Match(&req, &m, true); !b {
		t.Error(err)
	}
}

func TestPathVariables(t *testing.T) {
	req := mock.Request{}
	req.Path = "/a/b/c"

	m := mock.Definition{}
	m.Request.Path = "/a/:b/:c"
	m.Request.PathVariables = map[string]string{
	  "b": "[a-z]",
	  "c": "[a-z]",
	}

	mm := Request{}

	// match
	if b, err := mm.Match(&req, &m, true); !b {
		t.Error(err)
	}

	// not match
	req.Path = "a/b/3"
	
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}

	// not present
	req.Path = "a/b"
	
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}
}

func TestPathGlob(t *testing.T) {
	req := mock.Request{}
	req.Path = "/a/b/c"

	m := mock.Definition{}
	m.Request.Path = "/a/*"

	mm := Request{}

	if b, err := mm.Match(&req, &m, true); !b {
		t.Error(err)
	}
}

func TestMatchQueryString(t *testing.T) {
	rval := make(mock.Values)
	rval["test"] = []string{"test"}
	req := mock.Request{}
	req.QueryStringParameters = rval

	m := mock.Definition{}
	mval := make(mock.Values)
	mval["test"] = []string{"test"}
	m.Request.QueryStringParameters = mval

	mm := Request{}
	if b, err := mm.Match(&req, &m, true); !b {
		t.Error(err)
	}

	mval["test2"] = []string{"test2"}
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchQueryStringLenMismatch(t *testing.T) {
	rval := make(mock.Values)
	rval["test"] = []string{"test"}
	req := mock.Request{}
	req.QueryStringParameters = rval

	m := mock.Definition{}
	mval := make(mock.Values)
	mval["test"] = []string{"test", "test2"}
	m.Request.QueryStringParameters = mval

	mm := Request{}
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchQueryStringNonExisting(t *testing.T) {
	rval := make(mock.Values)
	rval["test2"] = []string{"test"}
	req := mock.Request{}
	req.QueryStringParameters = rval

	m := mock.Definition{}
	mval := make(mock.Values)
	mval["test"] = []string{"test"}
	m.Request.QueryStringParameters = mval

	mm := Request{}
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}
}
func TestMatchQueryStringGlob(t *testing.T) {
	rval := make(mock.Values)
	rval["test"] = []string{"test"}
	req := mock.Request{}
	req.QueryStringParameters = rval

	m := mock.Definition{}
	mval := make(mock.Values)
	mval["test"] = []string{"*es*"}
	m.Request.QueryStringParameters = mval

	mm := Request{}
	if b, err := mm.Match(&req, &m, true); !b {
		t.Error(err)
	}

	mval["test2"] = []string{"tes*"}
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchQueryMultiStringGlob(t *testing.T) {
	rval := make(mock.Values)
	rval["first"] = []string{"test"}
	rval["second"] = []string{"another_test"}
	req := mock.Request{}
	req.QueryStringParameters = rval

	m := mock.Definition{}
	mval := make(mock.Values)
	mval["first"] = []string{"t*"}
	mval["second"] = []string{"another_test"}
	m.Request.QueryStringParameters = mval

	mm := Request{}
	if b, err := mm.Match(&req, &m, true); !b {
		t.Error(err)
	}

	mval["first"] = []string{"t*"}
	mval["second"] = []string{"*another_test_2"}
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchQueryStringRegexp(t *testing.T) {
	rval := make(mock.Values)
	rval["test"] = []string{"test"}
	req := mock.Request{}
	req.QueryStringParameters = rval

	m := mock.Definition{}
	mval := make(mock.Values)
	mval["test"] = []string{"^t.*t$"}
	m.Request.QueryStringParameters = mval

	mm := Request{}
	if b, err := mm.Match(&req, &m, true); !b {
		t.Error(err)
	}

	mval["test2"] = []string{"tes.*\\d"}
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchQueryMultiStringRegexp(t *testing.T) {
	rval := make(mock.Values)
	rval["first"] = []string{"test"}
	rval["second"] = []string{"another_test"}
	req := mock.Request{}
	req.QueryStringParameters = rval

	m := mock.Definition{}
	mval := make(mock.Values)
	mval["first"] = []string{"t.*"}
	mval["second"] = []string{"another[-_@]test"}
	m.Request.QueryStringParameters = mval

	mm := Request{}
	if b, err := mm.Match(&req, &m, true); !b {
		t.Error(err)
	}

	mval["first"] = []string{"t*"}
	mval["second"] = []string{"*another_test_2"}
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}
}
func TestMatchCookies(t *testing.T) {
	rval := make(mock.Cookies)
	rval["test"] = "test"
	req := mock.Request{}
	req.Cookies = rval

	m := mock.Definition{}
	mval := make(mock.Cookies)
	mval["test"] = "test"
	m.Request.Cookies = mval

	mm := Request{}
	if b, err := mm.Match(&req, &m, true); !b {
		t.Error(err)
	}

	mval["test2"] = "test2"
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchCookiesNonExisting(t *testing.T) {
	rval := make(mock.Cookies)
	rval["test2"] = "test"
	req := mock.Request{}
	req.Cookies = rval

	m := mock.Definition{}
	mval := make(mock.Cookies)
	mval["test"] = "test"
	m.Request.Cookies = mval

	mm := Request{}
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchCookiesGlob(t *testing.T) {
	rval := make(mock.Cookies)
	rval["test"] = "test"
	req := mock.Request{}
	req.Cookies = rval

	m := mock.Definition{}
	mval := make(mock.Cookies)
	mval["test"] = "*es*"
	m.Request.Cookies = mval

	mm := Request{}
	if b, err := mm.Match(&req, &m, true); !b {
		t.Error(err)
	}

	mval["test2"] = "test*"
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchHeaders(t *testing.T) {
	rval := make(mock.Values)
	rval["test"] = []string{"test"}
	req := mock.Request{}
	req.Headers = rval

	m := mock.Definition{}
	mval := make(mock.Values)
	mval["test"] = []string{"test"}
	m.Request.Headers = mval

	mm := Request{}
	if b, err := mm.Match(&req, &m, true); !b {
		t.Error(err)
	}

	mval["test2"] = []string{"test2"}
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchHeadersGlobValues(t *testing.T) {
	rval := make(mock.Values)
	rval["test"] = []string{"test"}
	req := mock.Request{}
	req.Headers = rval

	m := mock.Definition{}
	mval := make(mock.Values)
	mval["test"] = []string{"*es*"}
	m.Request.Headers = mval

	mm := Request{}
	if b, err := mm.Match(&req, &m, true); !b {
		t.Error(err)
	}

	mval["test2"] = []string{"test*"}
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchHeadersMultiGlobValue(t *testing.T) {
	rval := make(mock.Values)
	rval["first"] = []string{"test"}
	rval["second"] = []string{"another_test"}
	req := mock.Request{}
	req.Headers = rval

	m := mock.Definition{}
	mval := make(mock.Values)
	mval["first"] = []string{"*es*"}
	mval["second"] = []string{"*ther_tes*"}
	m.Request.Headers = mval

	mm := Request{}
	if b, err := mm.Match(&req, &m, true); !b {
		t.Error(err)
	}

	mval["first"] = []string{"*es*"}
	mval["second"] = []string{"*tmher_es*"}
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchGlobHeadersGlobValue(t *testing.T) {
	rval := make(mock.Values)
	rval["test"] = []string{"test"}
	req := mock.Request{}
	req.Headers = rval

	m := mock.Definition{}
	mval := make(mock.Values)
	mval["*es*"] = []string{"*"}
	m.Request.Headers = mval

	mm := Request{}
	if b, err := mm.Match(&req, &m, true); !b {
		t.Error(err)
	}

	mval["wrong_one"] = []string{"invalid*"}
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchGlobHeadersGlobValueNonExisting(t *testing.T) {
	rval := make(mock.Values)
	rval["test"] = []string{"test"}
	req := mock.Request{}
	req.Headers = rval

	m := mock.Definition{}
	mval := make(mock.Values)
	mval["*invalid"] = []string{"*"}
	m.Request.Headers = mval

	mm := Request{}
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}

	mval = make(mock.Values)
	mval["*es*"] = []string{"*invalid"}
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchHost(t *testing.T) {

	req := mock.Request{}
	req.Host = "domain.com"

	m := mock.Definition{}
	m.Request.Host = "domain.com"

	mm := Request{}
	if b, err := mm.Match(&req, &m, true); !b {
		t.Error(err)
	}

	req.Host = "error.com"
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchHostGlob(t *testing.T) {

	req := mock.Request{}
	req.Host = "domain.com"

	m := mock.Definition{}
	m.Request.Host = "*omain.co*"

	mm := Request{}
	if b, err := mm.Match(&req, &m, true); !b {
		t.Error(err)
	}

	req.Host = "error.com"
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchBody(t *testing.T) {

	req := mock.Request{}
	req.Body = "HelloWorld"

	m := mock.Definition{}
	m.Request.Body = "HelloWorld"

	mm := Request{}
	if b, err := mm.Match(&req, &m, true); !b {
		t.Error(err)
	}

	req.Body = "ByeBye"
	if b, err := mm.Match(&req, &m, true); b {
		t.Error(err)
	}
}

func TestGlobBody(t *testing.T) {
	req := mock.Request{}
	req.Body = "Hello World From Test"

	m := mock.Definition{}
	m.Request.Body = "*World*"

	mm := Request{}
	if b, err := mm.Match(&req, &m, true); !b {
		t.Error(err)
	}

}

func TestBodyComparator(t *testing.T) {
	req := mock.Request{}
	req.Body = "{\"name\":\"bob\",\"age\":30}"
	hval := make(mock.Values)
	hval["Content-Type"] = []string{"application/json; charset=utf8"}
	req.Headers = hval

	m := mock.Definition{}
	m.Request.Body = "{\"age\":30,\n\"name\":\"bob\"}"

	comparator := payload.NewDefaultComparator()
	mm := Request{comparator: comparator}
	if b, err := mm.Match(&req, &m, true); !b {
		t.Error(err)
	}

}

func TestMatchIgnoreMissingBodyDefinition(t *testing.T) {
	req := mock.Request{}
	req.Body = "HelloWorld"
	m := mock.Definition{}
	mm := Request{}
	if b, _ := mm.Match(&req, &m, true); !b {
		t.Error("Not expected match")
	}
}

func TestSceneMatchingDefinition(t *testing.T) {
	req := mock.Request{}
	req.Body = "HelloWorld"
	m := mock.Definition{}
	m.Control.Scenario.Name = "uSEr"
	m.Control.Scenario.RequiredState = []string{"created"}
	s := NewInMemoryScenarioStore()
	mm := Request{scenario: s}
	if b, _ := mm.Match(&req, &m, true); b {
		t.Error("Scenario doesn't match")
	}
	s.SetState("user", "created")
	if b, _ := mm.Match(&req, &m, true); !b {
		t.Error("Scenario match")
	}
}

func TestSceneMatchingIgnoreStateCase(t *testing.T) {
	req := mock.Request{}
	req.Body = "HelloWorld"
	m := mock.Definition{}
	m.Control.Scenario.Name = "uSEr"
	m.Control.Scenario.RequiredState = []string{"CreAted"}
	s := NewInMemoryScenarioStore()
	mm := Request{scenario: s}
	if b, _ := mm.Match(&req, &m, true); b {
		t.Error("Scenario doesn't match")
	}
	s.SetState("user", "created")
	if b, _ := mm.Match(&req, &m, true); !b {
		t.Error("Scenario match")
	}
}

func TestSceneMatchingDefinitionDisabled(t *testing.T) {
	req := mock.Request{}
	req.Body = "HelloWorld"
	m := mock.Definition{}
	m.Control.Scenario.Name = "uSEr"
	m.Control.Scenario.RequiredState = []string{"created"}
	s := NewInMemoryScenarioStore()
	mm := Request{scenario: s}
	if b, _ := mm.Match(&req, &m, false); !b {
		t.Error("Scenario not skiped")
	}

	if b, _ := mm.Match(&req, &m, true); b {
		t.Error("Scenario skiped")
	}

}

func TestMatchIgnoreUnexpectedHeadersAndQuery(t *testing.T) {
	req := mock.Request{}
	req.Method = "GET"
	req.Path = "/a/b/c"
	hval := make(mock.Values)
	hval["test"] = []string{"test"}
	hval["test2"] = []string{"test"}
	hval["test3"] = []string{"test"}
	req.QueryStringParameters = hval
	req.Headers = hval

	m := mock.Definition{}
	m.Request.Method = "GET"
	m.Request.Path = "/a/b/c"
	mval := make(mock.Values)
	mval["test"] = []string{"test"}
	m.Request.QueryStringParameters = mval
	m.Request.Headers = mval

	mm := Request{}

	if b, _ := mm.Match(&req, &m, true); !b {
		t.Error("Not expected match")
	}
}
