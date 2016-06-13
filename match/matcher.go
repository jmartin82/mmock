package match

import (
	"mmock/definition"
)

type Matcher interface {
	Match(req *definition.Request, mock *definition.Request) (bool, error)
}
