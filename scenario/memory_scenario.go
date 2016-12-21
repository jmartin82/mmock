package scenario

import (
	"strings"
)

func NewInMemmoryScenarion() *InMemoryScenario {
	status := make(map[string]string)
	return &InMemoryScenario{status: status}
}

type InMemoryScenario struct {
	status map[string]string
}

func (sm *InMemoryScenario) SetState(name, status string) {
	sm.status[strings.ToLower(name)] = strings.ToLower(status)

}

func (sm *InMemoryScenario) GetState(name string) string {
	if v, f := sm.status[strings.ToLower(name)]; f {
		return v
	}
	return "not_started"
}
