package definition

//DefinitionReader interface contains the funtions to obtain the mock defintions.
type DefinitionReader interface {
	//ReadMockDefinition return an array of mock definitions sort by priority.
	ReadMocksDefinition() []Mock
}
