package scenario

type Director interface {
	SetState(name, status string)
	GetState(name string) string
	Reset(name string) bool
	ResetAll()
}
