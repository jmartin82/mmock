package mock

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"testing"
)

func TestBuildRequestDefinitionFromHTTP(t *testing.T) {
	b := bytes.NewBufferString("body text")
	req, _ := http.NewRequest("POST", "https://domain.tld:99901/test.php?aa=bb#fragment", b)
	req.TLS = &tls.ConnectionState{}
	req.Header.Add("X-TEST-HEADER", "random value")
	cookie := http.Cookie{Name: "cookie_name", Value: "cookie_value"}
	req.AddCookie(&cookie)

	tr := HTTP{}
	def := tr.BuildRequestDefinitionFromHTTP(req)

	if def.Scheme != "https" {
		t.Fatalf("Invalid scheme")
	}

	if def.Port != "99901" {
		t.Fatalf("Invalid Port")
	}

	if def.Method != "POST" {
		t.Fatalf("Invalid Method")
	}

	if c, f := def.Cookies["cookie_name"]; !f || c != "cookie_value" {
		t.Fatalf("Invalid Cookies")
	}

	if h, f := def.Headers["X-Test-Header"]; !f || h[0] != "random value" {
		t.Fatalf("Invalid Headers %v", def.Headers)
	}

	if q, f := def.QueryStringParameters["aa"]; !f || q[0] != "bb" {
		t.Fatalf("Invalid Query")
	}

	if def.Body != "body text" {
		t.Fatalf("Invalid Body: %v", def.Body)
	}

	if def.Path != "/test.php" {
		t.Fatalf("Invalid Path %v", def.Path)
	}

	if def.Host != "domain.tld" {
		t.Fatalf("Invalid Host %v", def.Host)
	}

	if def.Fragment != "fragment" {
		t.Fatalf("Invalid Fragment %v", def.Fragment)
	}

}
