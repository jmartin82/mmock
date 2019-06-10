package match

import (
	"sync"

	"github.com/jmartin82/mmock/pkg/mock"
)

type TransactionStorer interface {
	Save(Log)
	Reset()
	ResetMatch(mock.Request)
	GetAll() []Log
	Get(limit uint, offset uint) []Log
}


//InMemoryTransactionStore stores all received request and their matches in memory until the last reset
type InMemoryTransactionStore struct {
	matches []Log
	sync.Mutex
	checker Matcher
}

//Save store a match information
func (mrs *InMemoryTransactionStore) Save(req Log) {
	mrs.Lock()
	mrs.matches = append(mrs.matches, req)
	mrs.Unlock()
}

//Reset clean the request stored in memory
func (mrs *InMemoryTransactionStore) Reset() {
	mrs.Lock()
	mrs.matches = make([]Log, 0, 100)
	mrs.Unlock()
}

//ResetMatch clean the request stored in memory that matches a particular criteria
func (mrs *InMemoryTransactionStore) ResetMatch(req mock.Request) {
	matches := mrs.GetAll()
	mrs.Lock()
	var r = []Log{}
	for _, e := range matches {
		if c, _ := mrs.checker.Match(e.Request, &mock.Definition{Request: req}, false); c == false {
			r = append(r, e)
		}
	}

	mrs.matches = r
	mrs.Unlock()
}

//GetAll return current matches (positive and negative) in memory
func (mrs *InMemoryTransactionStore) GetAll() []Log {
	mrs.Lock()
	r := make([]Log, len(mrs.matches))
	copy(r, mrs.matches)
	mrs.Unlock()
	return r
}

//Get return an subset of current matches (positive and negative) in memory
func (mrs *InMemoryTransactionStore) Get(limit uint, offset uint) []Log {
	mrs.Lock()
	defer mrs.Unlock()

	max := offset + limit
	if max > uint(len(mrs.matches)) {
		max = uint(len(mrs.matches))
	}

	if offset >= max {
		return []Log{}
	}

	r := make([]Log, max-offset)
	copy(r, mrs.matches[offset:max])

	return r
}

//NewInMemoryScenarioStore is the InMemoryTransactionStore constructor
func NewInMemoryTransactionStore(checker Matcher) *InMemoryTransactionStore {
	reqs := make([]Log, 0, 100)
	return &InMemoryTransactionStore{matches: reqs, checker: checker}

}
