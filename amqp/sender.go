package amqp

import "github.com/jmartin82/mmock/definition"

//Sender sends messages to AMQP server
type Sender interface {
	//Send sends to amqp
	Send(per *definition.Persist, req *definition.Request, res *definition.Response) bool
}
