package match

import (
	"encoding/json"
	"strings"
)

type ScenearioStorer interface {
	SetState(name, status string)
	GetState(name string) string
	SetStateValues(name string, values map[string]string)
	SetStateValue(name, valueName, value string)
	GetStateValues(name string) map[string]string
	GetStateValue(name, valueName string) (string, bool)
	Reset(name string) bool
	ResetAll()
	SetPaused(newstate bool)
	GetPaused() bool
	List() string
}

type ScenearioStore struct {
}

func NewInMemoryScenarioStore() *InMemoryScenarioStore {
	status := make(map[string]string)
	values := make(map[string]map[string]string)
	return &InMemoryScenarioStore{
		status: status,
		values: values,
	}
}

type InMemoryScenarioStore struct {
	status map[string]string
	values map[string]map[string]string
	paused bool
	*ScenearioStore
}

func (sm *InMemoryScenarioStore) List() string {
	json, _ := json.MarshalIndent(sm.status, "", "  ")
	return string(json)
}

func (sm *InMemoryScenarioStore) Reset(name string) bool {
	if _, f := sm.status[strings.ToLower(name)]; f {
		sm.status[strings.ToLower(name)] = "not_started"
		sm.values[strings.ToLower(name)] = make(map[string]string)
		return true
	}
	return false
}

func (sm *InMemoryScenarioStore) ResetAll() {
	sm.status = make(map[string]string)
	sm.values = make(map[string](map[string]string))
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

func (sm *InMemoryScenarioStore) SetStateValues(name string, values map[string]string) {
	if sm.paused {
		return
	}
	sm.values[strings.ToLower(name)] = values
}

func (sm *InMemoryScenarioStore) SetStateValue(name, valueName, value string) {
	if sm.paused {
		return
	}
	sm.values[strings.ToLower(name)][strings.ToLower(valueName)] = value
}

func (sm *InMemoryScenarioStore) GetStateValues(name string) map[string]string {
	if v, f := sm.values[strings.ToLower(name)]; f {
		return v
	}

	return make(map[string]string)
}

func (sm *InMemoryScenarioStore) GetStateValue(name, valueName string) (string, bool) {
	if v, f := sm.values[strings.ToLower(name)][strings.ToLower(valueName)]; f {
		return v, true
	}

	return "", false
}

func (sm *InMemoryScenarioStore) GetPaused() bool {
	return sm.paused
}

func (sm *InMemoryScenarioStore) SetPaused(newstate bool) {
	sm.paused = newstate
}
