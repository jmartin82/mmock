package definition

import "os"

//ConfigReader interface allows recongnize if there is available some config reader for an a specific file.
type ConfigReader interface {
	CanRead(fileInfo os.FileInfo) bool
	Read(filename string) (Mock, error)
}
