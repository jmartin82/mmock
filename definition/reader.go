package definition

//Reader interface contains the funtions to obtain the mock defintions.
type Reader interface {
	//ReadMockDefinition return an array of mock definitions sort by priority.
	ReadMocksDefinition() []Mock
}
