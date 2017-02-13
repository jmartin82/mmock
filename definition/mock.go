package definition

type Scenario struct {
	Name          string   `json:"name"`
	RequiredState []string `json:"requiredState"`
	NewState      string   `json:"newState"`
}

type Control struct {
	Priority     int      `json:"priority"`
	Delay        int      `json:"delay"`
	Crazy        bool     `json:"crazy"`
	Scenario     Scenario `json:"scenario"`
	ProxyBaseURL string   `json:"proxyBaseURL"`
}

//Mock contains the user mock definition
type Mock struct {
	Name        string
	Description string   `json:"description"`
	Request     Request  `json:"request"`
	Response    Response `json:"response"`
	Control     Control  `json:"control"`
}
