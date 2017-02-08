package match

import (
	"github.com/jmartin82/mmock/definition"
)

type Spier interface {
	Find(definition.Request) []definition.Match
	GetMatched() []definition.Match
	GetNotMatched() []definition.Match
	Store
}
