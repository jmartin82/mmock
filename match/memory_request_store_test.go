package match

import (
	"testing"

	"github.com/jmartin82/mmock/definition"
)

func TestStoreRequest(t *testing.T) {

	msr := NewMemoryRequestStore()
	m1 := definition.Request{Host: "TEST1"}
	msr.Save(m1)
	m2 := definition.Request{Host: "TEST2"}
	msr.Save(m2)

	if len(msr.requests) != 2 {
		t.Fatalf("Invalid store len: %v", len(msr.requests))
	}

	if cap(msr.requests) != 100 {
		t.Fatalf("Invalid store cap: %v", cap(msr.requests))
	}

	if msr.requests[0].Host != "TEST1" || msr.requests[1].Host != "TEST2" {
		t.Fatalf("Invalid store content")
	}

}

func TestGetRequests(t *testing.T) {

	msr := NewMemoryRequestStore()
	m1 := definition.Request{Host: "TEST1"}
	msr.Save(m1)
	m2 := definition.Request{Host: "TEST2"}
	msr.Save(m2)

	reqs := msr.GetRequests()
	msr.Reset()

	if len(reqs) != 2 {
		t.Fatalf("Invalid store len: %v", len(reqs))
	}
	reqs = msr.GetRequests()

	if len(reqs) != 0 {
		t.Fatalf("Invalid store len: %v", len(reqs))
	}
}

func TestReset(t *testing.T) {

	msr := NewMemoryRequestStore()
	m1 := definition.Request{Host: "TEST1"}
	msr.Save(m1)
	m2 := definition.Request{Host: "TEST2"}
	msr.Save(m2)

	if len(msr.requests) != 2 {
		t.Fatalf("Invalid store len: %v", len(msr.requests))
	}

	if cap(msr.requests) != 100 {
		t.Fatalf("Invalid store cap: %v", cap(msr.requests))
	}

	msr.Reset()

	if len(msr.requests) != 0 {
		t.Fatalf("Invalid store len: %v", len(msr.requests))
	}

	if cap(msr.requests) != 100 {
		t.Fatalf("Invalid store cap: %v", cap(msr.requests))
	}

}
