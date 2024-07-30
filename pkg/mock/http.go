package mock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// HTTP is and adaptor beteewn the http and mock config.
type HTTP struct {
}

// BuildRequestDefinitionFromHTTP Read the request config and return a mock request.
func (t HTTP) BuildRequestDefinitionFromHTTP(req *http.Request) Request {

	res := Request{}
	res.Scheme = getScheme(req)
	res.Host, res.Port = getHostAndPort(req)
	res.Method = req.Method
	res.Path = req.URL.Path
	res.Fragment = req.URL.Fragment

	res.Headers = make(Values)
	for header, values := range req.Header {
		if header != "Cookie" {
			res.Headers[header] = values
		}
	}

	res.Cookies = make(Cookies)
	for _, cookie := range req.Cookies() {
		res.Cookies[cookie.Name] = cookie.Value
	}

	res.QueryStringParameters = make(Values)
	for name, values := range req.URL.Query() {
		res.QueryStringParameters[name] = values
	}

	body, _ := io.ReadAll(req.Body)
	res.Body = string(body)

	return res
}

func getScheme(req *http.Request) string {
	if req.TLS != nil {
		return "https"
	}

	return "http"
}

func getHostAndPort(req *http.Request) (string, string) {
	host := req.Host
	if len(host) == 0 {
		return "localhost", "80"
	}

	index := strings.Index(host, ":")
	if index > -1 {
		return host[0:index], host[index+1:]
	}

	return host, "80"
}

// WriteHTTPResponseFromDefinition read a mock response and write a http response.
func (t HTTP) WriteHTTPResponseFromDefinition(fr *Response, w http.ResponseWriter, req *http.Request) {
	if isSSE(fr) {
		streamResponse(fr, w, req)
		return
	}
	addHeadersAndCookies(fr, w)
	w.WriteHeader(fr.StatusCode)
	io.WriteString(w, fr.Body)
}

// Check if the response is of type SSE
func isSSE(fr *Response) bool {
	values, ok := fr.Headers["content-type"]
	if ok {
		for _, value := range values {
			return strings.ToLower(value) == "text/event-stream"
		}
	}
	return false
}

func addHeadersAndCookies(fr *Response, w http.ResponseWriter) {
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
}

// streamResponse - stream response
func streamResponse(fr *Response, w http.ResponseWriter, req *http.Request) {
	addHeadersAndCookies(fr, w)

	for _, response := range gjson.Parse(fr.Body).Array() {
		time.Sleep(time.Second * 2)
		select {
		case <-req.Context().Done():
			return
		default:
			ba, _ := json.Marshal((response.Value()))
			fmt.Fprintf(w, "data: %s\n\n", string(ba))
			w.(http.Flusher).Flush()
		}
	}
}
