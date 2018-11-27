package definition

//MatchError contains the tested uri and the match error
type MatchError struct {
	URI    string `json:"uri"`
	Reason string `json:"reason"`
}

//MatchResult contains the match result and the failing matches with different mocks and the reason or the fail.
type MatchResult struct {
	Found  bool         `json:"match"`
	URI    string       `json:"uri"`
	Errors []MatchError `json:"errors"`
}

//Match contains the whole information about the request match. The http request, the final response received and the matching result.
type Match struct {
	Time     int64        `json:"time"`
	Request  *Request     `json:"request"`
	Response *Response    `json:"response"`
	Result   *MatchResult `json:"result"`
}
