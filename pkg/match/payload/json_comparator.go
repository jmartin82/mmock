package payload

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"log"
)

type JSONComparator struct {
}

func isArray(s string) bool {
	st := strings.TrimLeft(s, " ")
	return len(st) > 0 && st[0] == '['
}

func (jc *JSONComparator) doCompareJSONRegex(jsonWithPatterns, jsonWithValues string) bool {
	var patterns map[string]interface{}
	var values map[string]interface{}
	if err := json.Unmarshal([]byte(jsonWithPatterns), &patterns); err != nil {
  log.Printf("error in json patterns: %v", err)
		return false
	}
	if err := json.Unmarshal([]byte(jsonWithValues), &values); err != nil {
  log.Printf("error in json values: %v", err)
		return false
	}
	return match(patterns, values)
}

func (jc *JSONComparator) doCompareArrayRegex(jsonWithPatterns, jsonWithValues string) bool {
	var patterns []map[string]interface{}
	var values []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonWithPatterns), &patterns); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(jsonWithValues), &values); err != nil {
		return false
	}
	for i := 0; i < len(patterns); i++ {
		if !match(patterns[i], values[i]) {
			return false
		}
	}
	return true
}

func match(p, v map[string]interface{}) bool {
	for field, pattern := range p {
		value, exists := v[field]
		if !exists {
			return false
		}
		str, ok := pattern.(string)
		if !ok {
			return reflect.DeepEqual(pattern, value)
		}
		matched, err := regexp.MatchString(str, fmt.Sprint(value))
		if err != nil || !matched {
			return false
		}
	}
	return true
}

func (jc *JSONComparator) Compare(s1, s2 string) bool {

	if isArray(s1) != isArray(s2) {
		return false
	}

	if isArray(s1) || isArray(s2) {
		return jc.doCompareArrayRegex(s1, s2)
	}

	return jc.doCompareJSONRegex(s1, s2)
}
