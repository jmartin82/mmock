package scenario

import "testing"

func TestBasicSenarioManage(t *testing.T) {
	ms := NewInMemoryScenario()

	state := ms.GetState("scene1")
	if state != "not_started" {
		t.Errorf("Invalid initial state")
	}

	ms.SetState("scEne1", "SOME_STATE")
	state = ms.GetState("Scene1")
	if state != "some_state" {
		t.Errorf("Invalid initial state")
	}
}
