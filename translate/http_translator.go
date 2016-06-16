package translate

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jmartin82/mmock/definition"
)

type HTTPTranslator struct {
}

func (t HTTPTranslator) BuildRequestDefinitionFromHTTP(req *http.Request) definition.Request {

	res := definition.Request{}
	res.Method = req.Method
	res.Path = req.URL.Path
	res.Headers = make(definition.Values)
	for header, values := range req.Header {
		if header != "Cookie" {
			res.Headers[header] = values
		}
	}

	res.Cookies = make(definition.Cookies)
	for _, cookie := range req.Cookies() {
		res.Cookies[cookie.Name] = cookie.Value
	}

	res.QueryStringParameters = make(definition.Values)
	for name, values := range req.URL.Query() {
		res.QueryStringParameters[name] = values
	}

	body, _ := ioutil.ReadAll(req.Body)
	res.Body = string(body)
	return res
}

func (t HTTPTranslator) WriteHTTPResponseFromDefinition(fr *definition.Response, w http.ResponseWriter) {

	for header, values := range fr.Headers {
		for _, value := range values {
			w.Header().Add(header, value)
		}

	}
	if len(fr.Cookies) > 0 {
		cookies := []string{}
		for cookie, value := range fr.Cookies {
			cookies = append(cookies, fmt.Sprintf("%s=%s", cookie, value))
		}
		w.Header().Add("Set-Cookie", strings.Join(cookies, ";"))
	}

	w.WriteHeader(fr.StatusCode)
	io.WriteString(w, fr.Body)
}
