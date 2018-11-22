package match

import (
	"reflect"
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

func TestGet(t *testing.T) {

	msr := NewMemoryStore()

	matches := []definition.Match{
		{Time: 1},
		{Time: 2},
		{Time: 3},
		{Time: 4},
		{Time: 5},
	}
	for _, m := range matches {
		msr.Save(m)
	}

	tests := []struct {
		msg      string
		limit    uint
		offset   uint
		expected []definition.Match
	}{
		{"Zero limit and offset", 0, 0, []definition.Match{}},
		{"Zero limit and one offset", 0, 1, []definition.Match{}},
		{"Grab the first element", 1, 0, []definition.Match{{Time: 1}}},
		{"Grab second element", 1, 1, []definition.Match{{Time: 2}}},
		{"Grab first two elements", 2, 0, []definition.Match{{Time: 1}, {Time: 2}}},
		{"Grab the second and the third elements", 2, 1, []definition.Match{{Time: 2}, {Time: 3}}},
		{"Grab the last elements", 1, 4, []definition.Match{{Time: 5}}},
		{"Grab the last two elements", 2, 3, []definition.Match{{Time: 4}, {Time: 5}}},
		{"Out of bounds offset", 1, 5, []definition.Match{}},
		{"Out of bounds limit", 2, 4, []definition.Match{{Time: 5}}},
	}

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			r := msr.Get(tt.limit, tt.offset)
			if !reflect.DeepEqual(r, tt.expected) {
				t.Errorf("Wrong definitions: got %v, want %v", r, tt.expected)
			}
		})
	}
}

func TestGetOnEmptyStore(t *testing.T) {

	msr := NewMemoryStore()

	tests := []struct {
		msg      string
		limit    uint
		offset   uint
		expected []definition.Match
	}{
		{"Zero limit and offset", 0, 0, []definition.Match{}},
		{"Zero limit and one offset", 0, 1, []definition.Match{}},
		{"Out of bounds offset", 0, 1, []definition.Match{}},
		{"Out of bounds limit", 1, 0, []definition.Match{}},
		{"Out of bounds", 1, 1, []definition.Match{}},
	}

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			r := msr.Get(tt.limit, tt.offset)
			if !reflect.DeepEqual(r, tt.expected) {
				t.Errorf("Wrong definitions: got %v, want %v", r, tt.expected)
			}
		})
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
