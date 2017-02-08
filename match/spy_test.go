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
	msr := NewMemoryStore()
	m1 := definition.Match{Request: &definition.Request{Host: "TEST1"}}
	msr.Save(m1)
	m2 := definition.Match{Request: &definition.Request{Host: "TEST2"}}
	msr.Save(m2)
	m3 := definition.Match{Request: &definition.Request{Host: "TEST1"}}
	msr.Save(m3)

	matchVeryfier := NewSpy(NewTester(DummyScenarioManager{}), msr)

	matches := matchVeryfier.Find(definition.Request{Host: "TEST1"})

	if len(matches) != 2 {
		t.Fatalf("Expected matches 2 != %v", len(matches))
	}

	for _, match := range matches {
		if match.Request.Host != "TEST1" {
			t.Fatalf("Invalid match")
		}
	}

}
