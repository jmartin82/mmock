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
	Persisted Persisted `json:"persisted"`
	Body string `json:"body"`
	BodyAppend string `json:"BodyAppend"`
}

type Persisted struct {
	Name string `json:"name"`
	NotFound struct {
		StatusCode int `json:"statusCode"`
		Body string `json:"body"`
		BodyAppend string `json:"BodyAppend"`
	} `json:"notFound"`
	BodyAppend string `json:"BodyAppend"`
}

type Persist struct {
	Name string `json:"name"`
	Delete bool `json:"delete"`
}
