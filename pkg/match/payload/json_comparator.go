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

func (jc *JSONComparator) doCompareJSONRegexUnmarshaled(patterns, values map[string]interface{}) bool {
	var matches bool
	matches = jc.match(patterns, values)
	if !matches{
		log.Debugf("values: %v don't match: %v", values, patterns)
	}
	return matches
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
	return jc.doCompareJSONRegexUnmarshaled(patterns, values)
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
	return jc.doCompareArrayRegexUnmarshaled(patterns, values)
}

func (jc *JSONComparator) doCompareArrayRegexUnmarshaled(patterns, values []map[string]interface{}) bool {

	for i := 0; i < len(patterns); i++ {
		if !jc.match(patterns[i], values[i]) {
				log.Debugf("value %v doesn't match %v",
					values[i], patterns[i])
			return false
		}
	}
	return true
}

func (jc *JSONComparator) match(p, v map[string]interface{}) bool {
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
			var valueType reflect.Kind
			var patternType reflect.Kind
			var result bool

			valueType = reflect.ValueOf(value).Kind()
			patternType = reflect.ValueOf(pattern).Kind()

			if valueType == reflect.Map && patternType == reflect.Map {
					log.Debugf("recursing into map %v", field)

				result = jc.doCompareJSONRegexUnmarshaled(
					pattern.(map[string]interface{}),
					value.(map[string]interface{}))

				if !result {
					return false
				}
			} else if (valueType == reflect.Array || valueType == reflect.Slice) &&
				(patternType == reflect.Array || patternType == reflect.Slice) {

					log.Debugf("recursing into array %v", field)
				valueJsonBytes, err1 := json.Marshal(value)
				patternJsonBytes, err2 := json.Marshal(pattern)

				if err1 != nil || err2 != nil {
						log.Errorf("value %v raised %v and pattern %v raised %v",
							value, err1, pattern, err2)
					return false
				}

				result = jc.doCompareArrayRegex(
					string(patternJsonBytes), string(valueJsonBytes))

				if !result {
					return false
				}
			} else {
				var eql bool
				eql = reflect.DeepEqual(pattern, value)

				if !eql {
						log.Debugf("value %v doesn't DeepEqual %v", value, pattern)
				}

				return eql
			}
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
