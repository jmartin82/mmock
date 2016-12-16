package scenario

type ScenarioManager interface {
	SetState(name, status string)
	GetState(name string) string
}
