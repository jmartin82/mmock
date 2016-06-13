package definition

type Values map[string][]string

type Cookies map[string]string

type headers struct {
	Headers Values  `json:"headers"`
	Cookies Cookies `json:"cookies"`
}

type Request struct {
	Method                string `json:"method"`
	Path                  string `json:"path"`
	QueryStringParameters Values `json:"queryStringParameters"`
	headers
	Body string `json:"body"`
}

type Response struct {
	StatusCode int `json:"statusCode"`
	headers
	Body string `json:"body"`
}
