package parse

import (
	"mmock/definition"
)

type ResponseParser interface {
	Parse(*definition.Request, *definition.Response)
}
