package amqp

import "github.com/jmartin82/mmock/definition"

//Sender sends messages to AMQP server
type Sender interface {
	//Send sends to amqp
	Send(m *definition.Mock) bool
}
