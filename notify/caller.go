package notify

import "github.com/jmartin82/mmock/definition"

//Caller makes remote http requests
type Caller interface {
	//Call makes a remote http request
	Call(m definition.Request) bool
}
