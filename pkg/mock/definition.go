package mock

import (
	"encoding/json"
	"fmt"
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

type Scenario struct {
	Name          string   `json:"name"`
	RequiredState []string `json:"requiredState"`
	NewState      string   `json:"newState"`
}

type Delay string

func (d *Delay) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if f, ok := v.(float64); ok {
		*d = Delay(fmt.Sprintf("%ds", int(f)))
		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		fmt.Println("! DEPRECATION NOTICE:                                        !")
		fmt.Println("! Please use a time unit (m,s,ms) to define the delay value. !")
		fmt.Println("! Ex: \"delay\":\"1ms\" instead \"delay\":1                  !")
		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		return nil
	}

	*d = Delay(v.(string))
	return nil
}

func (d Delay) Duration() (time.Duration, error) {
	return time.ParseDuration(string(d))
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
	Control     Control  `json:"control"`
}
