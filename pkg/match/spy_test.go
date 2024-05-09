package match

import (
	"testing"

	"github.com/jmartin82/mmock/v3/pkg/match/payload"
	"github.com/jmartin82/mmock/v3/pkg/mock"
)

type DummyScenarioManager struct {
}

func (dsm DummyScenarioManager) SetState(name, status string) {

}
func (dsm DummyScenarioManager) Reset(name string) bool {
	return true
}
func (dsm DummyScenarioManager) ResetAll() {

}
func (dsm DummyScenarioManager) GetState(name string) string {
	return ""
}

func (dsm DummyScenarioManager) GetStateValues(name string) map[string]string {
	return make(map[string]string)
}

func (dsm DummyScenarioManager) GetStateValue(name, valueName string) string {
	return ""
}

func (dsm DummyScenarioManager) SetStateValues(name string, values map[string]string) {
	return
}

func (dsm DummyScenarioManager) SetStateValue(name, valueName, value string) {
	return
}

func (dsm DummyScenarioManager) GetPaused() bool {
	return false
}

func (dsm DummyScenarioManager) SetPaused(_ bool) {
}

func (dsm DummyScenarioManager) List() string {
	return ""
}

func TestFindMatches(t *testing.T) {
	spy := NewSpy(NewTester(payload.NewComparator(), DummyScenarioManager{}), NewInMemoryTransactionStore(DummyMatcher{}, 10))

	m1 := Transaction{Request: &mock.Request{Host: "TEST1"}}
	spy.Save(m1)
	m2 := Transaction{Request: &mock.Request{Host: "TEST2"}}
	spy.Save(m2)
	m3 := Transaction{Request: &mock.Request{Host: "TEST1"}}
	spy.Save(m3)

	matches := spy.Find(mock.Request{Host: "TEST1"})

	if len(matches) != 2 {
		t.Fatalf("Expected matches 2 != %v", len(matches))
	}

	for _, match := range matches {
		if match.Request.Host != "TEST1" {
			t.Fatalf("Invalid match")
		}
	}

}

func TestMatchByResult(t *testing.T) {
	spy := NewSpy(NewTester(payload.NewComparator(), DummyScenarioManager{}), NewInMemoryTransactionStore(DummyMatcher{}, 10))

	m1 := Transaction{Result: &Result{Found: true}}
	spy.Save(m1)
	m2 := Transaction{Result: &Result{Found: false}}
	spy.Save(m2)
	m3 := Transaction{Result: &Result{Found: true}}
	spy.Save(m3)

	matches := spy.GetAll()

	if len(matches) != 3 {
		t.Fatalf("Expected matches 3 != %v", len(matches))
	}

	matches = spy.GetMatched()

	if len(matches) != 2 {
		t.Fatalf("Expected matches 2 != %v", len(matches))
	}
	matches = spy.GetUnMatched()

	if len(matches) != 1 {
		t.Fatalf("Expected matches 1 != %v", len(matches))
	}

}
