package definition

//Mock contains the user mock definition
type Mock struct {
	Name string
	Request Request `json:"request"`
	Response Response `json:"response"`
	Persist Persist `json:"store"`
	Control  struct {
		Priority     int    `json:"priority"`
		Delay        int    `json:"delay"`
		Crazy        bool   `json:"crazy"`
		ProxyBaseURL string `json:"proxyBaseURL"`
	} `json:"control"`
}
