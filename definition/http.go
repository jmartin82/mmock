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
	BodyAddition string `json:"bodyAddition"`
}

type Persisted struct {
	FileName string `json:"fileName"`
	NotFound struct {
		StatusCode int `json:"statusCode"`
		Body string `json:"body"`
		BodyAddition string `json:"bodyAddition"`
	} `json:"notFound"`
	BodyAddition string `json:"bodyAddition"`
}

type Persist struct {
	FileName string `json:"fileName"`
	Delete bool `json:"delete"`
}
