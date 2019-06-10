package match

import (
	"errors"
	"github.com/jmartin82/mmock/pkg/match/payload"

	"reflect"
	"testing"

	"github.com/jmartin82/mmock/pkg/mock"

)

type DummyMatcher struct {
	OK bool
}

func (dm DummyMatcher) Match(req *mock.Request, mock *mock.Definition, scenarioAware bool) (bool, error) {
	if dm.OK {
		return true, nil
	}
	return false, errors.New("Random Error")
}

func TestStoreRequest(t *testing.T) {

	msr := NewInMemoryTransactionStore(DummyMatcher{})
	m1 := Log{Request: &mock.Request{Host: "TEST1"}}
	msr.Save(m1)
	m2 := Log{Request: &mock.Request{Host: "TEST2"}}
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

	msr := NewInMemoryTransactionStore(DummyMatcher{})
	m1 := Log{Request: &mock.Request{Host: "TEST1"}}
	msr.Save(m1)
	m2 := Log{Request: &mock.Request{Host: "TEST2"}}
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

	msr := NewInMemoryTransactionStore(DummyMatcher{})

	matches := []Log{
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
		expected []Log
	}{
		{"Zero limit and offset", 0, 0, []Log{}},
		{"Zero limit and one offset", 0, 1, []Log{}},
		{"Grab the first element", 1, 0, []Log{{Time: 1}}},
		{"Grab second element", 1, 1, []Log{{Time: 2}}},
		{"Grab first two elements", 2, 0, []Log{{Time: 1}, {Time: 2}}},
		{"Grab the second and the third elements", 2, 1, []Log{{Time: 2}, {Time: 3}}},
		{"Grab the last elements", 1, 4, []Log{{Time: 5}}},
		{"Grab the last two elements", 2, 3, []Log{{Time: 4}, {Time: 5}}},
		{"Out of bounds offset", 1, 5, []Log{}},
		{"Out of bounds limit", 2, 4, []Log{{Time: 5}}},
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

	msr := NewInMemoryTransactionStore(DummyMatcher{})

	tests := []struct {
		msg      string
		limit    uint
		offset   uint
		expected []Log
	}{
		{"Zero limit and offset", 0, 0, []Log{}},
		{"Zero limit and one offset", 0, 1, []Log{}},
		{"Out of bounds offset", 0, 1, []Log{}},
		{"Out of bounds limit", 1, 0, []Log{}},
		{"Out of bounds", 1, 1, []Log{}},
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

	msr := NewInMemoryTransactionStore(DummyMatcher{})
	m1 := Log{Request: &mock.Request{Host: "TEST1"}}
	msr.Save(m1)
	m2 := Log{Request: &mock.Request{Host: "TEST2"}}
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

func TestResetMatch(t *testing.T) {

	scenario := NewInMemoryScenarioStore()
	comparator := payload.NewDefaultComparator()
	tester := NewTester(comparator, scenario)

	msr := NewInMemoryTransactionStore(tester)
	m1 := Log{Request: &mock.Request{Host: "TEST1"}}
	msr.Save(m1)
	m2 := Log{Request: &mock.Request{Host: "TEST2"}}
	msr.Save(m2)

	if len(msr.matches) != 2 {
		t.Fatalf("Invalid store len: %v", len(msr.matches))
	}

	if cap(msr.matches) != 100 {
		t.Fatalf("Invalid store cap: %v", cap(msr.matches))
	}

	msr.ResetMatch(mock.Request{
		Host: "TEST1",
	})

	if len(msr.matches) != 1 {
		t.Fatalf("Invalid store len: %v", len(msr.matches))
	}

	if cap(msr.matches) != 1 {
		t.Fatalf("Invalid store cap: %v", cap(msr.matches))
	}

}
