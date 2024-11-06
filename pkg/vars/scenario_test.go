package vars

import (
	"strings"
	"testing"

	"github.com/jmartin82/mmock/v3/pkg/match"
)

var scenarioTestVars = []struct {
	key          string
	value        string
	expectToFind bool
}{
	{"scenario.test1", "one", true},
	{"scenario.test2", "two", true},
	{"scenario.not_found", "nothing", false},
}

func getLoadedScenarioFiller() ScenarioFiller {
	var scenarioName = "test_scenario"

	var store = match.NewInMemoryScenarioStore()
	store.SetStateValue(scenarioName, "test1", "one")
	store.SetStateValue(scenarioName, "test2", "two")

	filler := ScenarioFiller{
		Name:  scenarioName,
		Store: store,
	}

	return filler
}

func TestScenarioFiller(t *testing.T) {
	var filler = getLoadedScenarioFiller()

	for _, tt := range scenarioTestVars {
		holders := []string{
			tt.key,
		}

		vars := filler.Fill(holders)

		if len(vars) == 0 {
			if tt.expectToFind {
				t.Errorf("Unable to retrieve vars")
			}
			continue
		}

		v, ok := vars[tt.key]

		if !ok {
			t.Errorf("Unable to retrieve value inside vars")
			continue
		}

		if strings.EqualFold(v[0], tt.value) != tt.expectToFind {
			t.Errorf("Couldn't get the expected value. Expected: %s, Value found: %s", tt.value, v[0])
		}
	}
}
