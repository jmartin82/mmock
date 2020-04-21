package mock

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

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

type Response struct {
	StatusCode int `json:"statusCode"`
	HttpHeaders
	Body string `json:"body"`
}

type Callback struct {
	Delay       Delay  `json:"delay"`
	Url         string `json:"url"`
	ContentType string `json:"contentType"`
	Body        string `json:"body"`
}

type Scenario struct {
	Name          string   `json:"name"`
	RequiredState []string `json:"requiredState"`
	NewState      string   `json:"newState"`
}

type Delay struct {
	time.Duration
}

func (d *Delay) UnmarshalJSON(data []byte) (err error) {
	var (
		v interface{}
		s string
	)
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v.(type) {
	case float64:
		s = fmt.Sprintf("%ds", int(v.(float64)))
		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		fmt.Println("! DEPRECATION NOTICE:                                        !")
		fmt.Println("! Please use a time unit (m,s,ms) to define the delay value. !")
		fmt.Println("! Ex: \"delay\":\"1s\" instead \"delay\":1                   !")
		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	case string:
		s = v.(string)
	default:
		return fmt.Errorf("invalid value for delay, got: %v", reflect.TypeOf(v))
	}

	d.Duration, err = time.ParseDuration(s)
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

//Definition contains the user mock config
type Definition struct {
	URI         string
	Description string   `json:"description"`
	Request     Request  `json:"request"`
	Response    Response `json:"response"`
	Callback    Callback `json:"callback"`
	Control     Control  `json:"control"`
}
