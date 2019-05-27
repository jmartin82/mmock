package payload

import (
	"encoding/json"
	"reflect"
)

type JSONComparator struct {
}

func isArray(s1 string) bool {
	return len(s1) > 0 && s1[0] == '['
}

func (jc *JSONComparator) CompareArray(s1, s2 string) bool {

	var o1 []interface{}
	var o2 []interface{}

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

	if isArray(s1) || isArray(s2) {
		return jc.CompareArray(s1, s2)
	}

	var o1 interface{}
	var o2 interface{}

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
