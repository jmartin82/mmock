package match

import (
	"github.com/jmartin82/mmock/definition"
)

type Verifier interface {
	Verify(definition.Request) []definition.Request
}
