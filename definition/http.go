package definition

import "regexp"

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

//GetURLPart getting url part by passing regex pattern and groupName
func (req Request) GetURLPart(pattern string, groupName string) (string, bool) {
	r, error := regexp.Compile(pattern)
	if error != nil{
		return "", false
	} 

	match := r.FindStringSubmatch(req.Path)
	result := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			result[name] = match[i]
		}
	}

	value, present := result[groupName]

	return value, present
}

//GetQueryStringParam getting query string parameter value from request
func (req Request) GetQueryStringParam(name string) (string, bool) {

	if len(req.QueryStringParameters) == 0 {
		return "", false
	}
	value, f := req.QueryStringParameters[name]
	if !f {
		return "", false
	}

	return value[0], true
}

//GetCookieParam getting cookie param from request
func (req Request) GetCookieParam(name string) (string, bool) {

	if len(req.Cookies) == 0 {
		return "", false
	}
	value, f := req.Cookies[name]
	if !f {
		return "", false
	}

	return value, true
}