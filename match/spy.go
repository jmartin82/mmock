package match

import "github.com/jmartin82/mmock/definition"

type Spy struct {
	store   Store
	checker Checker
}

func NewSpy(checker Checker, matchStore Store) *Spy {
	return &Spy{store: matchStore, checker: checker}
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

// ResetMatch ...
func (mc Spy) ResetMatch(r definition.Request) {
	mc.store.ResetMatch(r)
}

func (mc Spy) Save(match definition.Match) {
	mc.store.Save(match)
}
func (mc Spy) Reset() {
	mc.store.Reset()
}

func (mc Spy) GetAll() []definition.Match {
	return mc.store.GetAll()
}

func (mc Spy) Get(limit uint, offset uint) []definition.Match {
	return mc.store.Get(limit, offset)
}

func (mc Spy) GetMatched() []definition.Match {
	return mc.getMatchByResult(true)
}

func (mc Spy) GetUnMatched() []definition.Match {
	return mc.getMatchByResult(false)
}

func (mc Spy) getMatchByResult(found bool) []definition.Match {
	matches := mc.store.GetAll()
	result := []definition.Match{}
	for _, match := range matches {
		if match.Result.Found == found {
			result = append(result, match)
		}
	}
	return result
}
