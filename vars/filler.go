package vars

import "github.com/jmartin82/mmock/definition"

type Filler interface {
	Fill(m *definition.Mock, input string) string
}
