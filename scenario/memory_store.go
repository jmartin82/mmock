package scenario

import (
	"strings"
)

func NewMemoryStore() *MemoryStore {
	status := make(map[string]string)
	return &MemoryStore{status: status}
}

type MemoryStore struct {
	status map[string]string
}

func (sm *MemoryStore) SetState(name, status string) {
	sm.status[strings.ToLower(name)] = strings.ToLower(status)
}

func (sm *MemoryStore) GetState(name string) string {
	if v, f := sm.status[strings.ToLower(name)]; f {
		return v
	}
	return "not_started"
}
