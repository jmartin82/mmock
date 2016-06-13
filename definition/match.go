package definition

type Result struct {
	Found  bool              `json:"match"`
	Errors map[string]string `json:"errors"`
}

type Match struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
	Result   Result   `json:"result"`
}
