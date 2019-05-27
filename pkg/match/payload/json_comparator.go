package payload

import (
	"encoding/json"
	"reflect"
	"strings"
)

type JSONComparator struct {
}

func isArray(s string) bool {
	st := strings.TrimLeft(s, " ")
	return len(st) > 0 && st[0] == '['
}

func (js *JSONComparator) doCompare(s1, s2 string, o1, o2 interface{}) bool {
	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false
	}

	return reflect.DeepEqual(o1, o2)
}

func (jc *JSONComparator) Compare(s1, s2 string) bool {

	if isArray(s1) != isArray(s2) {
		return false
	}

	if isArray(s1) || isArray(s2) {
		var o1 []interface{}
		var o2 []interface{}
		return jc.doCompare(s1, s2, o1, o2)
	}

	var o1 interface{}
	var o2 interface{}

	return jc.doCompare(s1, s2, o1, o2)
}
