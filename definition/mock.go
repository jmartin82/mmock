package definition

type Control struct {
	Priority     int    `json:"priority"`
	Delay        int    `json:"delay"`
	Crazy        bool   `json:"crazy"`
	ProxyBaseURL string `json:"proxyBaseURL"`
}

type Actions map[string]string

type Persist struct {
	Entity  string  `json:"entity"`
	Actions Actions `json:"actions"`
	Engine  string  `json:"engine"`
}

//Mock contains the user mock definition
type Mock struct {
	Name        string
	Description string   `json:"description"`
	Request     Request  `json:"request"`
	Response    Response `json:"response"`
	Persist     Persist  `json:"persist"`
	Control     Control  `json:"control"`
}
