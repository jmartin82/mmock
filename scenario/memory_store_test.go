package scenario

import "testing"

func TestBasicSenarioManage(t *testing.T) {
	ms := NewMemoryStore()

	state := ms.GetState("scene1")
	if state != "not_started" {
		t.Errorf("Invalid initial state")
	}

	ms.SetState("scEne1", "SOME_STATE")
	state = ms.GetState("Scene1")
	if state != "some_state" {
		t.Errorf("Invalid initial state")
	}

	ms.SetState("Scene1", "some_state")
	ms.SetState("Scene2", "some_state")

	if ms.GetState("Scene1") != "some_state" || ms.GetState("Scene2") != "some_state" {
		t.Errorf("Invalid initial state")
	}

	ms.Reset("Scene1")
	if ms.GetState("Scene1") != "not_started" || ms.GetState("Scene2") != "some_state" {
		t.Errorf("Invalid initial state")
	}

	ms.ResetAll()
	if ms.GetState("Scene1") != "not_started" || ms.GetState("Scene2") != "not_started" {
		t.Errorf("Invalid initial state")
	}

}
