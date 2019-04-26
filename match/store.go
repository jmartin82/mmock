package match

import (
	"github.com/jmartin82/mmock/definition"
)

type Store interface {
	Save(definition.Match)
	Reset()
	ResetMatch(definition.Request)
	GetAll() []definition.Match
	Get(limit uint, offset uint) []definition.Match
}
