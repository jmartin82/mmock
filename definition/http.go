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
	BodyAddition string `json:"bodyAddition"`
}

type File struct {
	Name    string `json:"name"`
	Delete  bool   `json:"delete"`
	Content struct {
		Name     string
		Request  Request  `json:"request"`
		Response Response `json:"response"`
		Control  Control  `json:"control"`
	} `json:"content"`
}

type Control struct {
	Priority     int    `json:"priority"`
	Delay        int    `json:"delay"`
	Crazy        bool   `json:"crazy"`
	ProxyBaseURL string `json:"proxyBaseURL"`
}
