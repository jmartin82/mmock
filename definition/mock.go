package definition

type Control struct {
	Priority     int    `json:"priority"`
	Delay        int    `json:"delay"`
	Crazy        bool   `json:"crazy"`
	ProxyBaseURL string `json:"proxyBaseURL"`
}

type Actions map[string]string
type Requests []Request

type Persist struct {
	Entity     string  `json:"entity"`
	Collection string  `json:"collection"`
	Actions    Actions `json:"actions"`
	Engine     string  `json:"engine"`
}

type Notify struct {
	Amqp AMQPPublishing `json:"amqp"`
	Http Requests       `json:"http"`
}

//Mock contains the user mock definition
type Mock struct {
	Name        string
	Description string   `json:"description"`
	Request     Request  `json:"request"`
	Response    Response `json:"response"`
	Persist     Persist  `json:"persist"`
	Notify      Notify   `json:"notify"`
	Control     Control  `json:"control"`
}
