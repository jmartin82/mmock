package mock

import (
	"encoding/json"
	"fmt"
	"github.com/jmartin82/mmock/v3/internal/config/logger"
	"reflect"
	"time"
)

var log = logger.Log

type Values map[string][]string

type Cookies map[string]string

type HttpHeaders struct {
	Headers Values  `json:"headers"`
	Cookies Cookies `json:"cookies"`
}

type Request struct {
	Scheme                string `json:"scheme"`
	Host                  string `json:"host"`
	Port                  string `json:"port"`
	Method                string `json:"method"`
	Path                  string `json:"path"`
	QueryStringParameters Values `json:"queryStringParameters"`
	Fragment              string `json:"fragment"`
	HttpHeaders
	Body string `json:"body"`
}

type ReplacementRequiredPayload interface {
	GetHeaders() HttpHeaders
	GetBody() string
}

type ReplacementRequired struct {
}

func (rr Response) GetHeaders() HttpHeaders {
	return rr.HttpHeaders
}

func (rr Response) GetBody() string {
	return rr.Body
}

func (rr Callback) GetHeaders() HttpHeaders {
	return rr.HttpHeaders
}

func (rr Callback) GetBody() string {
	return rr.Body
}

type Response struct {
	StatusCode int `json:"statusCode"`
	HttpHeaders
	Body string `json:"body"`
	*ReplacementRequired
}

type Callback struct {
	Delay  Delay  `json:"delay"`
	Method string `json:"method"`
	Url    string `json:"url"`
	HttpHeaders
	Body    string `json:"body"`
	Timeout Delay  `json:"timeout"`
	*ReplacementRequired
}

type Scenario struct {
	Name          string   `json:"name"`
	RequiredState []string `json:"requiredState"`
	NewState      string   `json:"newState"`
}

type Delay struct {
	time.Duration
}

func (d Delay) MarshalJSON() (data []byte, err error) {
	return json.Marshal(d.Duration.String())
}

func (d *Delay) UnmarshalJSON(data []byte) (err error) {

	var i interface{}
	if err = json.Unmarshal(data, &i); err != nil {
		return err
	}

	switch v := i.(type) {
	case string:
		d.Duration, err = time.ParseDuration(v)
	case time.Duration:
		d.Duration = v
	default:
		return fmt.Errorf("invalid value for delay got: %v. Ex: \"delay\":\"1s\" ", reflect.TypeOf(v))
	}

	return err
}

type Control struct {
	Priority     int      `json:"priority"`
	Delay        Delay    `json:"delay"`
	Crazy        bool     `json:"crazy"`
	Scenario     Scenario `json:"scenario"`
	ProxyBaseURL string   `json:"proxyBaseURL"`
	WebHookURL   string   `json:"webHookURL"`
}

// Definition contains the user mock config
type Definition struct {
	URI         string
	Description string   `json:"description"`
	Request     Request  `json:"request"`
	Response    Response `json:"response"`
	Callback    Callback `json:"callback"`
	Control     Control  `json:"control"`
}

func CreateDefinition(request Request, response Response) Definition {
	definition := Definition{Request: request, Response: response}
	log.Debugf("definition: %v", definition)
	return definition
}
