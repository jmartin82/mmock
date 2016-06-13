package match

import (
	"github.com/jmartin82/mmock/definition"
)

type Matcher interface {
	Match(req *definition.Request, mock *definition.Request) (bool, error)
}
