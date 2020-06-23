package match

import (
	"github.com/jmartin82/mmock/v3/pkg/mock"
)

type TransactionSpier interface {
	Find(mock.Request) []Transaction
	GetMatched() []Transaction
	GetUnMatched() []Transaction
	TransactionStorer
}

type Spy struct {
	store   TransactionStorer
	checker Matcher
}

func NewSpy(checker Matcher, transactionStore TransactionStorer) *Spy {
	return &Spy{store: transactionStore, checker: checker}
}

func (mc Spy) Find(r mock.Request) []Transaction {
	matches := mc.store.GetAll()
	result := []Transaction{}
	for _, match := range matches {
		if m, _ := mc.checker.Match(match.Request, &mock.Definition{Request: r}, false); m {
			result = append(result, match)
		}
	}
	return result

}

// ResetMatch ...
func (mc Spy) ResetMatch(r mock.Request) {
	mc.store.ResetMatch(r)
}

func (mc Spy) Save(match Transaction) {
	mc.store.Save(match)
}
func (mc Spy) Reset() {
	mc.store.Reset()
}

func (mc Spy) GetAll() []Transaction {
	return mc.store.GetAll()
}

func (mc Spy) Get(limit int, offset int) []Transaction {
	return mc.store.Get(limit, offset)
}

func (mc Spy) GetMatched() []Transaction {
	return mc.getMatchByResult(true)
}

func (mc Spy) GetUnMatched() []Transaction {
	return mc.getMatchByResult(false)
}

func (mc Spy) getMatchByResult(found bool) []Transaction {
	matches := mc.store.GetAll()
	result := []Transaction{}
	for _, match := range matches {
		if match.Result.Found == found {
			result = append(result, match)
		}
	}
	return result
}
