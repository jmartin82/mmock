package definition

//Result contains the match result and the failing matches with diferent mocks and the reason or the fail.
type Result struct {
	Found  bool              `json:"match"`
	Errors map[string]string `json:"errors"`
}

//Match contains the whole information about the request match. The http request, the final response recieved and the matching result.
type Match struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
	Result   Result   `json:"result"`
}
