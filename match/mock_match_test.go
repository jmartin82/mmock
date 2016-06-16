package match

import (
	"testing"

	"github.com/jmartin82/mmock/definition"
)

func TestMatchMethod(t *testing.T) {
	hreq := &definition.Request{}
	hreq.Method = "GET"
	mreq := &definition.Request{}
	mreq.Method = "GET"
	m := MockMatch{}

	if m, err := m.Match(hreq, mreq); !m {
		t.Error(err)
	}

	mreq.Method = "POST"
	if m, _ := m.Match(hreq, mreq); m {
		t.Error("Not expected match")
	}
}

func TestMatchPath(t *testing.T) {
	hreq := &definition.Request{}
	hreq.Path = "/a/b/c"
	mreq := &definition.Request{}
	mreq.Path = "/a/b/c"
	m := MockMatch{}

	if m, err := m.Match(hreq, mreq); !m {
		t.Error(err)
	}

	mreq.Path = "/a/b/d"
	if m, _ := m.Match(hreq, mreq); m {
		t.Error("Not expected match")
	}
}

func TestGlobPath(t *testing.T) {
	hreq := &definition.Request{}
	hreq.Path = "/a/b/c"
	mreq := &definition.Request{}
	mreq.Path = "/a/b/*"
	m := MockMatch{}

	if m, err := m.Match(hreq, mreq); !m {
		t.Error(err)
	}
}

func TestMatchQueryString(t *testing.T) {

	hreq := &definition.Request{}
	hval := make(definition.Values)
	hval["test"] = []string{"test"}
	hreq.QueryStringParameters = hval

	mreq := &definition.Request{}
	mval := make(definition.Values)
	mval["test"] = []string{"test"}
	mreq.QueryStringParameters = mval

	m := MockMatch{}

	if m, err := m.Match(hreq, mreq); !m {
		t.Error(err)
	}

	mval["test2"] = []string{"test2"}
	if m, _ := m.Match(hreq, mreq); m {
		t.Error("Not expected match")
	}

}

func TestMatchCookies(t *testing.T) {
	hreq := &definition.Request{}
	hval := make(definition.Cookies)
	hval["cookie"] = "val"
	hreq.Cookies = hval

	mreq := &definition.Request{}
	mval := make(definition.Cookies)
	mval["cookie"] = "val"
	mreq.Cookies = mval

	m := MockMatch{}

	if m, err := m.Match(hreq, mreq); !m {
		t.Error(err)
	}

	mval["cookie2"] = "val2"
	if m, _ := m.Match(hreq, mreq); m {
		t.Error("Not expected match")
	}
}

func TestMatchHeaders(t *testing.T) {
	hreq := &definition.Request{}
	hval := make(definition.Values)
	hval["test"] = []string{"test"}
	hreq.Headers = hval

	mreq := &definition.Request{}
	mval := make(definition.Values)
	mval["test"] = []string{"test"}
	mreq.Headers = mval

	m := MockMatch{}

	if m, err := m.Match(hreq, mreq); !m {
		t.Error(err)
	}

	mval["test"] = []string{"test2"}
	if m, _ := m.Match(hreq, mreq); m {
		t.Error("Not expected match")
	}
}

func TestMatchBody(t *testing.T) {
	hreq := &definition.Request{}
	hreq.Body = "HelloWorld"
	mreq := &definition.Request{}
	mreq.Body = "HelloWorld"
	m := MockMatch{}

	if m, err := m.Match(hreq, mreq); !m {
		t.Error(err)
	}

	mreq.Body = "ByeWorld"
	if m, _ := m.Match(hreq, mreq); m {
		t.Error("Not expected match")
	}
}

func TestGlobBody(t *testing.T) {
	hreq := &definition.Request{}
	hreq.Body = "Hello World From Test"
	mreq := &definition.Request{}
	mreq.Body = "*World*"
	m := MockMatch{}

	if m, err := m.Match(hreq, mreq); !m {
		t.Error(err)
	}

}

func TestMatchIgnoreUnexpectedHeadersAnQuery(t *testing.T) {
	hreq := &definition.Request{}
	hreq.Method = "GET"
	hreq.Path = "/a/b/c"
	hval := make(definition.Values)
	hval["test"] = []string{"test"}
	hval["test2"] = []string{"test"}
	hval["test3"] = []string{"test"}
	hreq.QueryStringParameters = hval
	hreq.Headers = hval

	mreq := &definition.Request{}
	mreq.Method = "GET"
	mreq.Path = "/a/b/c"
	mval := make(definition.Values)
	mval["test"] = []string{"test"}
	mreq.QueryStringParameters = hval
	mreq.Headers = mval

	m := MockMatch{}

	if m, err := m.Match(hreq, mreq); !m {
		t.Error(err)
	}
}
