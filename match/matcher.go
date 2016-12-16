package match

import (
	"github.com/jmartin82/mmock/definition"
)

//Matcher checks if the received request matches with some specific mock request definition.
type Matcher interface {
	Match(req *definition.Request, mock *definition.Mock) (bool, error)
}
