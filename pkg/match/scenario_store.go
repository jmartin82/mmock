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


func NewScenarioStore() *ScenarioStore {
	status := make(map[string]string)
	return &ScenarioStore{status: status}
}

type ScenarioStore struct {
	status map[string]string
	paused bool
}

func (sm *ScenarioStore) Reset(name string) bool {
	if _, f := sm.status[strings.ToLower(name)]; f {
		sm.status[strings.ToLower(name)] = "not_started"
		return true
	}
	return false
}

func (sm *ScenarioStore) ResetAll() {
	sm.status = make(map[string]string)
	sm.paused = false
}

func (sm *ScenarioStore) SetState(name, status string) {
	if sm.paused {
		return
	}
	sm.status[strings.ToLower(name)] = strings.ToLower(status)
}

func (sm *ScenarioStore) GetState(name string) string {
	if v, f := sm.status[strings.ToLower(name)]; f {
		return v
	}
	return "not_started"
}

func (sm *ScenarioStore) GetPaused() bool {
	return sm.paused
}

func (sm *ScenarioStore) SetPaused(newstate bool) {
	sm.paused = newstate
}
