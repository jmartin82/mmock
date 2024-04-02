package match

import (
	"time"

	"github.com/jmartin82/mmock/v3/pkg/mock"
)

//Error contains the tested uri and the match error
type Error struct {
	URI    string `json:"uri"`
	Reason string `json:"reason"`
}

//Result contains the match result and the failing matches with different mocks and the reason or the fail.
type Result struct {
	Found  bool    `json:"match"`
	URI    string  `json:"uri"`
	Errors []Error `json:"errors"`
}

//Transaction contains the whole information about the request match. The http request, the final response received and the matching result.
type Transaction struct {
	Time     int64          `json:"time"`
	Request  *mock.Request  `json:"request"`
	Response *mock.Response `json:"response"`
	Result   *Result        `json:"result"`
}

func NewTransaction(request *mock.Request, response *mock.Response, result *Result) *Transaction {

	return &Transaction{Time: time.Now().Unix(), Request: request, Response: response, Result: result}
}
