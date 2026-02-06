package vars

import (
	"github.com/jmartin82/mmock/v3/pkg/mock"
)

type DummyScenarioStorer struct {
	Name   string
	Values mock.ScenarioValues
}

func NewDummyScenarioStorer(name string, values mock.ScenarioValues) DummyScenarioStorer {
	result := DummyScenarioStorer{
		Name:   name,
		Values: values,
	}

	return result
}

func (dss DummyScenarioStorer) SetState(name, status string) {
	return
}

func (dss DummyScenarioStorer) GetState(name string) string {
	return ""
}

func (dss DummyScenarioStorer) SetStateValues(name string, values map[string]string) {
	return

}

func (dss DummyScenarioStorer) SetStateValue(name, valueName, value string) {
	return
}

func (dss DummyScenarioStorer) GetStateValues(name string) map[string]string {
	return make(map[string]string)
}

func (dss DummyScenarioStorer) GetStateValue(name, valueName string) (string, bool) {
	return "", false
}

func (dss DummyScenarioStorer) Reset(name string) bool {
	return true
}

func (dss DummyScenarioStorer) ResetAll() {
	return
}

func (dss DummyScenarioStorer) SetPaused(newstate bool) {
	return
}

func (dss DummyScenarioStorer) GetPaused() bool {
	return false
}

func (dss DummyScenarioStorer) List() string {
	return ""
}
