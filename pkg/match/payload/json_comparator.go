package payload

import (
	"encoding/json"
	"fmt"
	"github.com/jmartin82/mmock/v3/internal/config/logger"
	"reflect"
	"regexp"
	"strings"
)

var log = logger.Log

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
		log.Errorf("error in json patterns: %v", err)
		return false
	}

	if err := json.Unmarshal([]byte(jsonWithValues), &values); err != nil {
		log.Errorf("error in json values: %v", err)
		return false
	}
	return match(patterns, values)
}

func (jc *JSONComparator) doCompareArrayRegex(jsonWithPatterns, jsonWithValues string) bool {
	var patterns []map[string]interface{}
	var values []map[string]interface{}

	if err := json.Unmarshal([]byte(jsonWithPatterns), &patterns); err != nil {
		log.Errorf("error in json patterns: %v", err)

		return false
	}
	if err := json.Unmarshal([]byte(jsonWithValues), &values); err != nil {
		log.Errorf("error in json patterns: %v", err)
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
		log.Debugf("comparing field %v with pattern %v against value %v",
			field, pattern, value)

		if !exists {
			log.Debugf("field doesn't exist: %v", field)

			return false
		}
		str, ok := pattern.(string)
		if !ok {
			return reflect.DeepEqual(pattern, value)
		}
		matched, err := regexp.MatchString(str, fmt.Sprint(value))
		if err != nil || !matched {
			log.Debugf("value %v doesn't match %v : %v", fmt.Sprint(value), str, err)

			return false
		}
	}
	return true
}

func (jc *JSONComparator) Compare(s1, s2 string) bool {

	if isArray(s1) != isArray(s2) {
		log.Debugf("only one of these is an array %v %v", s1, s2)

		return false
	}

	if isArray(s1) || isArray(s2) {
		return jc.doCompareArrayRegex(s1, s2)
	}

	return jc.doCompareJSONRegex(s1, s2)
}
