package definition

//Mock contains the user mock definition
type Mock struct {
	Name     string
	Request  Request  `json:"request"`
	Response Response `json:"response"`
	File     File     `json:"file"`
	Control  Control  `json:"control"`
}
