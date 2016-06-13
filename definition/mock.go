package definition

type Mock struct {
	Name     string
	Request  Request  `json:"request"`
	Response Response `json:"response"`
	Control  struct {
		Priority int  `json:"priority"`
		Delay    int  `json:"delay"`
		Crazy    bool `json:"crazy"`
	} `json:"control"`
}
