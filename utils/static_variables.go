package utils

// Variables is used for storing static variables
type StaticVariables struct {
	ServerAddress string
}

var variables = StaticVariables{}

// SetServerAddress sets ServerAddress variable
func SetServerAddress(address string) {
	variables.ServerAddress = address
}

// GetServerAddress returns ServerAddress variable
func GetServerAddress() string {
	return variables.ServerAddress
}
