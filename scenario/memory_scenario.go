package scenario

func NewInMemmoryScenarion() *InMemoryScenario {
	status := make(map[string]string)
	return &InMemoryScenario{status: status}
}

type InMemoryScenario struct {
	status map[string]string
}

func (sm *InMemoryScenario) SetState(name, status string) {
	sm.status[name] = status

}

func (sm *InMemoryScenario) GetState(name string) string {
	if v, f := sm.status[name]; f {
		return v
	}
	return "not_started"
}
