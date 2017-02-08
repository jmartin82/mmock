package definition

type MatchErrors map[string]string

//Result contains the match result and the failing matches with different mocks and the reason or the fail.
type Result struct {
	Found  bool        `json:"match"`
	Errors MatchErrors `json:"errors"`
}

//Match contains the whole information about the request match. The http request, the final response received and the matching result.
type Match struct {
	Time     int64     `json:"time"`
	Request  *Request  `json:"request"`
	Response *Response `json:"response"`
	Result   Result    `json:"result"`
}
