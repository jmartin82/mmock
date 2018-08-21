package match

import (
	"testing"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/scenario"
)

func TestMatchMethod(t *testing.T) {
	req := definition.Request{}
	req.Method = "GET"

	m := definition.Mock{}
	m.Request.Method = "GET"

	mm := Tester{}
	if b, err := mm.Check(&req, &m, true); !b {
		t.Error(err)
	}

	req.Method = "POST"
	if b, err := mm.Check(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchScheme(t *testing.T) {
	req := definition.Request{}
	req.Scheme = "https"

	m := definition.Mock{}
	m.Request.Scheme = "https"

	mm := Tester{}
	if b, err := mm.Check(&req, &m, true); !b {
		t.Error(err)
	}

	req.Scheme = "http"
	if b, err := mm.Check(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchFragment(t *testing.T) {
	req := definition.Request{}
	req.Fragment = "fragment"

	m := definition.Mock{}
	m.Request.Fragment = "fragment"

	mm := Tester{}
	if b, err := mm.Check(&req, &m, true); !b {
		t.Error(err)
	}

	req.Fragment = "nothing"
	if b, err := mm.Check(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchPath(t *testing.T) {

	req := definition.Request{}
	req.Path = "/a/b/c"

	m := definition.Mock{}
	m.Request.Path = "/a/b/c"

	mm := Tester{}
	if b, err := mm.Check(&req, &m, true); !b {
		t.Error(err)
	}

	req.Path = "/a/b/d"
	if b, err := mm.Check(&req, &m, true); b {
		t.Error(err)
	}
}

func TestPathVars(t *testing.T) {
	req := definition.Request{}
	req.Path = "/a/b/c"

	m := definition.Mock{}
	m.Request.Path = "/a/:b/:c"

	mm := Tester{}

	if b, err := mm.Check(&req, &m, true); !b {
		t.Error(err)
	}
}

func TestPathGlob(t *testing.T) {
	req := definition.Request{}
	req.Path = "/a/b/c"

	m := definition.Mock{}
	m.Request.Path = "/a/*"

	mm := Tester{}

	if b, err := mm.Check(&req, &m, true); !b {
		t.Error(err)
	}
}

func TestMatchQueryString(t *testing.T) {
	rval := make(definition.Values)
	rval["test"] = []string{"test"}
	req := definition.Request{}
	req.QueryStringParameters = rval

	m := definition.Mock{}
	mval := make(definition.Values)
	mval["test"] = []string{"test"}
	m.Request.QueryStringParameters = mval

	mm := Tester{}
	if b, err := mm.Check(&req, &m, true); !b {
		t.Error(err)
	}

	mval["test2"] = []string{"test2"}
	if b, err := mm.Check(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchQueryStringLenMismatch(t *testing.T) {
	rval := make(definition.Values)
	rval["test"] = []string{"test"}
	req := definition.Request{}
	req.QueryStringParameters = rval

	m := definition.Mock{}
	mval := make(definition.Values)
	mval["test"] = []string{"test", "test2"}
	m.Request.QueryStringParameters = mval

	mm := Tester{}
	if b, err := mm.Check(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchQueryStringNonExisting(t *testing.T) {
	rval := make(definition.Values)
	rval["test2"] = []string{"test"}
	req := definition.Request{}
	req.QueryStringParameters = rval

	m := definition.Mock{}
	mval := make(definition.Values)
	mval["test"] = []string{"test"}
	m.Request.QueryStringParameters = mval

	mm := Tester{}
	if b, err := mm.Check(&req, &m, true); b {
		t.Error(err)
	}
}
func TestMatchQueryStringGlob(t *testing.T) {
	rval := make(definition.Values)
	rval["test"] = []string{"test"}
	req := definition.Request{}
	req.QueryStringParameters = rval

	m := definition.Mock{}
	mval := make(definition.Values)
	mval["test"] = []string{"*es*"}
	m.Request.QueryStringParameters = mval

	mm := Tester{}
	if b, err := mm.Check(&req, &m, true); !b {
		t.Error(err)
	}

	mval["test2"] = []string{"tes*"}
	if b, err := mm.Check(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchQueryMultiStringGlob(t *testing.T) {
	rval := make(definition.Values)
	rval["first"] = []string{"test"}
	rval["second"] = []string{"another_test"}
	req := definition.Request{}
	req.QueryStringParameters = rval

	m := definition.Mock{}
	mval := make(definition.Values)
	mval["first"] = []string{"t*"}
	mval["second"] = []string{"another_test"}
	m.Request.QueryStringParameters = mval

	mm := Tester{}
	if b, err := mm.Check(&req, &m, true); !b {
		t.Error(err)
	}

	mval["first"] = []string{"t*"}
	mval["second"] = []string{"*another_test_2"}
	if b, err := mm.Check(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchCookies(t *testing.T) {
	rval := make(definition.Cookies)
	rval["test"] = "test"
	req := definition.Request{}
	req.Cookies = rval

	m := definition.Mock{}
	mval := make(definition.Cookies)
	mval["test"] = "test"
	m.Request.Cookies = mval

	mm := Tester{}
	if b, err := mm.Check(&req, &m, true); !b {
		t.Error(err)
	}

	mval["test2"] = "test2"
	if b, err := mm.Check(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchCookiesNonExisting(t *testing.T) {
	rval := make(definition.Cookies)
	rval["test2"] = "test"
	req := definition.Request{}
	req.Cookies = rval

	m := definition.Mock{}
	mval := make(definition.Cookies)
	mval["test"] = "test"
	m.Request.Cookies = mval

	mm := Tester{}
	if b, err := mm.Check(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchCookiesGlob(t *testing.T) {
	rval := make(definition.Cookies)
	rval["test"] = "test"
	req := definition.Request{}
	req.Cookies = rval

	m := definition.Mock{}
	mval := make(definition.Cookies)
	mval["test"] = "*es*"
	m.Request.Cookies = mval

	mm := Tester{}
	if b, err := mm.Check(&req, &m, true); !b {
		t.Error(err)
	}

	mval["test2"] = "test*"
	if b, err := mm.Check(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchHeaders(t *testing.T) {
	rval := make(definition.Values)
	rval["test"] = []string{"test"}
	req := definition.Request{}
	req.Headers = rval

	m := definition.Mock{}
	mval := make(definition.Values)
	mval["test"] = []string{"test"}
	m.Request.Headers = mval

	mm := Tester{}
	if b, err := mm.Check(&req, &m, true); !b {
		t.Error(err)
	}

	mval["test2"] = []string{"test2"}
	if b, err := mm.Check(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchHeadersGlobValues(t *testing.T) {
	rval := make(definition.Values)
	rval["test"] = []string{"test"}
	req := definition.Request{}
	req.Headers = rval

	m := definition.Mock{}
	mval := make(definition.Values)
	mval["test"] = []string{"*es*"}
	m.Request.Headers = mval

	mm := Tester{}
	if b, err := mm.Check(&req, &m, true); !b {
		t.Error(err)
	}

	mval["test2"] = []string{"test*"}
	if b, err := mm.Check(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchHeadersMultiGlobValue(t *testing.T) {
	rval := make(definition.Values)
	rval["first"] = []string{"test"}
	rval["second"] = []string{"another_test"}
	req := definition.Request{}
	req.Headers = rval

	m := definition.Mock{}
	mval := make(definition.Values)
	mval["first"] = []string{"*es*"}
	mval["second"] = []string{"*ther_tes*"}
	m.Request.Headers = mval

	mm := Tester{}
	if b, err := mm.Check(&req, &m, true); !b {
		t.Error(err)
	}

	mval["first"] = []string{"*es*"}
	mval["second"] = []string{"*tmher_es*"}
	if b, err := mm.Check(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchGlobHeadersGlobValue(t *testing.T) {
	rval := make(definition.Values)
	rval["test"] = []string{"test"}
	req := definition.Request{}
	req.Headers = rval

	m := definition.Mock{}
	mval := make(definition.Values)
	mval["*es*"] = []string{"*"}
	m.Request.Headers = mval

	mm := Tester{}
	if b, err := mm.Check(&req, &m, true); !b {
		t.Error(err)
	}

	mval["wrong_one"] = []string{"invalid*"}
	if b, err := mm.Check(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchGlobHeadersGlobValueNonExisting(t *testing.T) {
	rval := make(definition.Values)
	rval["test"] = []string{"test"}
	req := definition.Request{}
	req.Headers = rval

	m := definition.Mock{}
	mval := make(definition.Values)
	mval["*invalid"] = []string{"*"}
	m.Request.Headers = mval

	mm := Tester{}
	if b, err := mm.Check(&req, &m, true); b {
		t.Error(err)
	}

	mval = make(definition.Values)
	mval["*es*"] = []string{"*invalid"}
	if b, err := mm.Check(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchHost(t *testing.T) {

	req := definition.Request{}
	req.Host = "domain.com"

	m := definition.Mock{}
	m.Request.Host = "domain.com"

	mm := Tester{}
	if b, err := mm.Check(&req, &m, true); !b {
		t.Error(err)
	}

	req.Host = "error.com"
	if b, err := mm.Check(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchHostGlob(t *testing.T) {

	req := definition.Request{}
	req.Host = "domain.com"

	m := definition.Mock{}
	m.Request.Host = "*omain.co*"

	mm := Tester{}
	if b, err := mm.Check(&req, &m, true); !b {
		t.Error(err)
	}

	req.Host = "error.com"
	if b, err := mm.Check(&req, &m, true); b {
		t.Error(err)
	}
}

func TestMatchBody(t *testing.T) {

	req := definition.Request{}
	req.Body = "HelloWorld"

	m := definition.Mock{}
	m.Request.Body = "HelloWorld"

	mm := Tester{}
	if b, err := mm.Check(&req, &m, true); !b {
		t.Error(err)
	}

	req.Body = "ByeBye"
	if b, err := mm.Check(&req, &m, true); b {
		t.Error(err)
	}
}

func TestGlobBody(t *testing.T) {
	req := definition.Request{}
	req.Body = "Hello World From Test"

	m := definition.Mock{}
	m.Request.Body = "*World*"

	mm := Tester{}
	if b, err := mm.Check(&req, &m, true); !b {
		t.Error(err)
	}

}

func TestMatchIgnoreMissingBodyDefinition(t *testing.T) {
	req := definition.Request{}
	req.Body = "HelloWorld"
	m := definition.Mock{}
	mm := Tester{}
	if b, _ := mm.Check(&req, &m, true); !b {
		t.Error("Not expected match")
	}
}

func TestSceneMatchingDefinition(t *testing.T) {
	req := definition.Request{}
	req.Body = "HelloWorld"
	m := definition.Mock{}
	m.Control.Scenario.Name = "uSEr"
	m.Control.Scenario.RequiredState = []string{"created"}
	s := scenario.NewMemoryStore()
	mm := Tester{Scenario: s}
	if b, _ := mm.Check(&req, &m, true); b {
		t.Error("Scenario doesn't match")
	}
	s.SetState("user", "created")
	if b, _ := mm.Check(&req, &m, true); !b {
		t.Error("Scenario match")
	}
}

func TestSceneMatchingIgnoreStateCase(t *testing.T) {
	req := definition.Request{}
	req.Body = "HelloWorld"
	m := definition.Mock{}
	m.Control.Scenario.Name = "uSEr"
	m.Control.Scenario.RequiredState = []string{"CreAted"}
	s := scenario.NewMemoryStore()
	mm := Tester{Scenario: s}
	if b, _ := mm.Check(&req, &m, true); b {
		t.Error("Scenario doesn't match")
	}
	s.SetState("user", "created")
	if b, _ := mm.Check(&req, &m, true); !b {
		t.Error("Scenario match")
	}
}

func TestSceneMatchingDefinitionDisabled(t *testing.T) {
	req := definition.Request{}
	req.Body = "HelloWorld"
	m := definition.Mock{}
	m.Control.Scenario.Name = "uSEr"
	m.Control.Scenario.RequiredState = []string{"created"}
	s := scenario.NewMemoryStore()
	mm := Tester{Scenario: s}
	if b, _ := mm.Check(&req, &m, false); !b {
		t.Error("Scenario not skiped")
	}

	if b, _ := mm.Check(&req, &m, true); b {
		t.Error("Scenario skiped")
	}

}

func TestMatchIgnoreUnexpectedHeadersAndQuery(t *testing.T) {
	req := definition.Request{}
	req.Method = "GET"
	req.Path = "/a/b/c"
	hval := make(definition.Values)
	hval["test"] = []string{"test"}
	hval["test2"] = []string{"test"}
	hval["test3"] = []string{"test"}
	req.QueryStringParameters = hval
	req.Headers = hval

	m := definition.Mock{}
	m.Request.Method = "GET"
	m.Request.Path = "/a/b/c"
	mval := make(definition.Values)
	mval["test"] = []string{"test"}
	m.Request.QueryStringParameters = mval
	m.Request.Headers = mval

	mm := Tester{}

	if b, _ := mm.Check(&req, &m, true); !b {
		t.Error("Not expected match")
	}
}
