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
	Delay  Delay  `json:"delay"`
	Method string `json:"method"`
	Url    string `json:"url"`
	HttpHeaders
	Body    string `json:"body"`
	Timeout Delay  `json:"timeout"`
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
	case float64:
		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		fmt.Println("! DEPRECATION NOTICE:                                        !")
		fmt.Println("! Please use a time unit (m,s,ms) to define the delay value. !")
		fmt.Println("! Ex: \"delay\":\"1s\" instead \"delay\":1                   !")
		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		s := fmt.Sprintf("%ds", int(v))
		d.Duration, err = time.ParseDuration(s)
	case string:
		d.Duration, err = time.ParseDuration(v)
	case time.Duration:
		d.Duration = v
	default:
		return fmt.Errorf("invalid value for delay, got: %v", reflect.TypeOf(v))
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
