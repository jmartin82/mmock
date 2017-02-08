package match

import (
	"testing"

	"github.com/jmartin82/mmock/definition"
)

func TestStoreRequest(t *testing.T) {

	msr := NewMemoryStore()
	m1 := definition.Match{Request: &definition.Request{Host: "TEST1"}}
	msr.Save(m1)
	m2 := definition.Match{Request: &definition.Request{Host: "TEST2"}}
	msr.Save(m2)

	if len(msr.matches) != 2 {
		t.Fatalf("Invalid store len: %v", len(msr.matches))
	}

	if cap(msr.matches) != 100 {
		t.Fatalf("Invalid store cap: %v", cap(msr.matches))
	}

	if msr.matches[0].Request.Host != "TEST1" || msr.matches[1].Request.Host != "TEST2" {
		t.Fatalf("Invalid store content")
	}

}

func TestGetAll(t *testing.T) {

	msr := NewMemoryStore()
	m1 := definition.Match{Request: &definition.Request{Host: "TEST1"}}
	msr.Save(m1)
	m2 := definition.Match{Request: &definition.Request{Host: "TEST2"}}
	msr.Save(m2)

	reqs := msr.GetAll()
	msr.Reset()

	if len(reqs) != 2 {
		t.Fatalf("Invalid store len: %v", len(reqs))
	}
	reqs = msr.GetAll()

	if len(reqs) != 0 {
		t.Fatalf("Invalid store len: %v", len(reqs))
	}
}

func TestReset(t *testing.T) {

	msr := NewMemoryStore()
	m1 := definition.Match{Request: &definition.Request{Host: "TEST1"}}
	msr.Save(m1)
	m2 := definition.Match{Request: &definition.Request{Host: "TEST2"}}
	msr.Save(m2)

	if len(msr.matches) != 2 {
		t.Fatalf("Invalid store len: %v", len(msr.matches))
	}

	if cap(msr.matches) != 100 {
		t.Fatalf("Invalid store cap: %v", cap(msr.matches))
	}

	msr.Reset()

	if len(msr.matches) != 0 {
		t.Fatalf("Invalid store len: %v", len(msr.matches))
	}

	if cap(msr.matches) != 100 {
		t.Fatalf("Invalid store cap: %v", cap(msr.matches))
	}

}
