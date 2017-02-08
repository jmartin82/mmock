package match

import "github.com/jmartin82/mmock/definition"

type Spy struct {
	store   Store
	checker Checker
}

func (mc Spy) Find(r definition.Request) []definition.Match {
	matches := mc.store.GetAll()
	result := []definition.Match{}
	for _, match := range matches {
		if m, _ := mc.checker.Check(match.Request, &definition.Mock{Request: r}, false); m {
			result = append(result, match)
		}
	}
	return result

}

func (mc Spy) GetAll() []definition.Match {
	return mc.store.GetAll()
}

func (mc Spy) GetMatched() []definition.Match {
	return []definition.Match{}
}

func (mc Spy) GetNotMatched() []definition.Match {
	return []definition.Match{}
}
func (mc Spy) Forget() {
	mc.store.Reset()
}

func NewSpy(checker Checker, matchStore Store) *Spy {
	return &Spy{store: matchStore, checker: checker}
}
