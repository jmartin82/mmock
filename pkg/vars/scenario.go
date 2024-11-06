package vars

import (
	"github.com/jmartin82/mmock/v3/pkg/match"
	"strings"
)

type ScenarioFiller struct {
	Name  string
	Store match.ScenearioStorer
}

func (sf ScenarioFiller) Fill(holders []string) map[string][]string {
	vars := make(map[string][]string)
	for _, tag := range holders {
		found := false
		s := ""

		if strings.HasPrefix(tag, "scenario.") {
			s, found = sf.Store.GetStateValue(sf.Name, tag[9:])
		}

		if found {
			vars[tag] = append(vars[tag], s)
		}

	}
	return vars
}
