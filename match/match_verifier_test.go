package match

import (
	"testing"

	"github.com/jmartin82/mmock/definition"
)

type DummyScenarioManager struct {
}

func (dsm DummyScenarioManager) SetState(name, status string) {

}
func (dsm DummyScenarioManager) GetState(name string) string {
	return ""
}

func TestMach(t *testing.T) {
	msr := NewMemoryRequestStore()
	m1 := definition.Request{Host: "TEST1"}
	msr.Save(m1)
	m2 := definition.Request{Host: "TEST2"}
	msr.Save(m2)
	m3 := definition.Request{Host: "TEST1"}
	msr.Save(m3)

	matchVeryfier := NewMatchVerifier(MockMatcher{Scenario: DummyScenarioManager{}}, msr)

	matches := matchVeryfier.Verify(definition.Request{Host: "TEST1"})

	if len(matches) != 2 {
		t.Fatalf("Expected matches 2 != %v", len(matches))
	}

	for _, match := range matches {
		if match.Host != "TEST1" {
			t.Fatalf("Invalid match")
		}
	}

}
