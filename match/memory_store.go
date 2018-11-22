package match

import (
	"sync"

	"github.com/jmartin82/mmock/definition"
)

//MemoryStore stores all received request and their matches in memory until the last reset
type MemoryStore struct {
	matches []definition.Match
	sync.Mutex
}

//Save store a match information
func (mrs *MemoryStore) Save(req definition.Match) {
	mrs.Lock()
	mrs.matches = append(mrs.matches, req)
	mrs.Unlock()
}

//Reset clean the request stored in memory
func (mrs *MemoryStore) Reset() {
	mrs.Lock()
	mrs.matches = make([]definition.Match, 0, 100)
	mrs.Unlock()
}

//GetAll return current matches (positive and negative) in memory
func (mrs *MemoryStore) GetAll() []definition.Match {
	mrs.Lock()
	r := make([]definition.Match, len(mrs.matches))
	copy(r, mrs.matches)
	mrs.Unlock()
	return r
}

//Get return an subset of current matches (positive and negative) in memory
func (mrs *MemoryStore) Get(limit uint, offset uint) []definition.Match {
	mrs.Lock()
	defer mrs.Unlock()

	max := offset + limit
	if max > uint(len(mrs.matches)) {
		max = uint(len(mrs.matches))
	}

	if offset >= max {
		return []definition.Match{}
	}

	r := make([]definition.Match, max-offset)
	copy(r, mrs.matches[offset:max])

	return r
}

//NewMemoryStore is the MemoryStore contructor
func NewMemoryStore() *MemoryStore {
	reqs := make([]definition.Match, 0, 100)
	return &MemoryStore{matches: reqs}

}
