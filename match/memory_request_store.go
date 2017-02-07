package match

import (
	"github.com/jmartin82/mmock/definition"
)

//MemoryRequestStore stores all received request in memory until the last reset
type MemoryRequestStore struct {
	requests []definition.Request
}

//Save store a new request on memory
func (mrs *MemoryRequestStore) Save(req definition.Request) {

	mrs.requests = append(mrs.requests, req)
}

//Reset clean the request stored in memory
func (mrs *MemoryRequestStore) Reset() {
	mrs.requests = make([]definition.Request, 0, 100)
}

//GetRequests return current requests in memory
func (mrs *MemoryRequestStore) GetRequests() []definition.Request {
	r := make([]definition.Request, len(mrs.requests))
	copy(r, mrs.requests)
	return r
}

//NewMemoryRequestStore is the MemoryRequestStore contructor
func NewMemoryRequestStore() *MemoryRequestStore {
	reqs := make([]definition.Request, 0, 100)
	return &MemoryRequestStore{requests: reqs}

}
