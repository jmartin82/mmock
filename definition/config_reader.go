package definition

//ConfigReader interface allows recognize if there is available some config reader for an a specific file.
type ConfigReader interface {
	CanRead(filename string) bool
	Read(content []byte) (Mock, error)
}
