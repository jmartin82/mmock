package definition

type DefinitionReader interface {
	ReadMocksDefinition() []Mock
}
