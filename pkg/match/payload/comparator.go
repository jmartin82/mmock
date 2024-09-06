package payload

import (
	"strings"
)

type Comparer interface {
	Compare(s1, s2 string, optionalPaths map[string]bool, currentPath string) bool
}

type Comparator struct {
	comparers map[string]Comparer
}

func NewComparator() *Comparator {
	comparers := make(map[string]Comparer)
	return &Comparator{comparers: comparers}
}

func NewDefaultComparator() *Comparator {
	comparator := NewComparator()
	json := &JSONComparator{}
	xml := &XMLComparator{}
	comparator.AddComparer("application/json", json)
	comparator.AddComparer("application/ld+json", json)
	comparator.AddComparer("application/merge-patch+json", json)
	comparator.AddComparer("application/xml", xml)
	comparator.AddComparer("text/xml", xml)
	return comparator
}

func (c Comparator) AddComparer(contentType string, comparer Comparer) {
	c.comparers[contentType] = comparer
}

func (c Comparator) Compare(contentType, s1, s2 string, optionalPaths map[string]bool) (comparable bool, equals bool) {
	parts := strings.Split(contentType, ";")
	comparer, ok := c.comparers[parts[0]]
	if !ok {
		return false, false
	}

	return true, comparer.Compare(s1, s2, optionalPaths, "")
}
