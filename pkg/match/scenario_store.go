package match

import (
	"strings"
)

type ScenearioStorer interface {
	SetState(name, status string)
	GetState(name string) string
	Reset(name string) bool
	ResetAll()
	SetPaused(newstate bool)
	GetPaused() bool
}

func NewInMemoryScenarioStore() *InMemoryScenarioStore {
	status := make(map[string]string)
	return &InMemoryScenarioStore{status: status}
}

type InMemoryScenarioStore struct {
	status map[string]string
	paused bool
}

func (sm *InMemoryScenarioStore) Reset(name string) bool {
	if _, f := sm.status[strings.ToLower(name)]; f {
		sm.status[strings.ToLower(name)] = "not_started"
		return true
	}
	return false
}

func (sm *InMemoryScenarioStore) ResetAll() {
	sm.status = make(map[string]string)
	sm.paused = false
}

func (sm *InMemoryScenarioStore) SetState(name, status string) {
	if sm.paused {
		return
	}
	sm.status[strings.ToLower(name)] = strings.ToLower(status)
}

func (sm *InMemoryScenarioStore) GetState(name string) string {
	if v, f := sm.status[strings.ToLower(name)]; f {
		return v
	}
	return "not_started"
}

func (sm *InMemoryScenarioStore) GetPaused() bool {
	return sm.paused
}

func (sm *InMemoryScenarioStore) SetPaused(newstate bool) {
	sm.paused = newstate
}
