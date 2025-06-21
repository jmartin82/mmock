package payload

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/jmartin82/mmock/v3/internal/config/logger"
)

var log = logger.Log

type JSONComparator struct {
}

func isArray(s string) bool {
	st := strings.TrimLeft(s, " ")
	return len(st) > 0 && st[0] == '['
}

func (jc *JSONComparator) doCompareJSONRegexUnmarshaled(
	patterns, values map[string]interface{},
	optionalPaths map[string]bool,
	currentPath string) bool {
	var matches bool
	matches = jc.match(patterns, values, optionalPaths, currentPath)
	if !matches {
		log.Debugf("values: %v don't match: %v", values, patterns)
	}
	return matches
}

func (jc *JSONComparator) doCompareJSONRegex(
	jsonWithPatterns, jsonWithValues string,
	optionalPaths map[string]bool,
	currentPath string) bool {
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
	return jc.doCompareJSONRegexUnmarshaled(patterns, values, optionalPaths, currentPath)
}

func (jc *JSONComparator) doCompareArrayRegex(
	jsonWithPatterns, jsonWithValues string,
	optionalPaths map[string]bool,
	currentPath string) bool {
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
	return jc.doCompareArrayRegexUnmarshaled(patterns, values, optionalPaths, currentPath)
}

func (jc *JSONComparator) doCompareArrayRegexUnmarshaled(
	patterns, values []map[string]interface{},
	optionalPaths map[string]bool,
	currentPath string) bool {

	for i := 0; i < len(patterns); i++ {
		if !jc.match(patterns[i], values[i], optionalPaths, currentPath) {
			log.Debugf("value %v doesn't match %v",
				values[i], patterns[i])
			return false
		}
	}
	return true
}

func (jc *JSONComparator) match(
	p, v map[string]interface{},
	optionalPaths map[string]bool,
	currentPath string) bool {
	for field, pattern := range p {
		var currentFieldPath = fmt.Sprintf("%s.%s", currentPath, field)
		value, exists := v[field]
		log.Debugf("comparing field %v with pattern %v against value %v",
			currentFieldPath, pattern, value)

		if !exists && !optionalPaths[currentFieldPath] {
			log.Debugf("field doesn't exist, and isn't optional: %v", currentFieldPath)

			return false
		} else if !exists && optionalPaths[currentFieldPath] {

			log.Debugf("field doesn't exist, but is optional: %v", currentFieldPath)
			continue
		}

		str, ok := pattern.(string)
		if !ok {
			var valueType reflect.Kind
			var patternType reflect.Kind
			var result bool

			valueType = reflect.ValueOf(value).Kind()
			patternType = reflect.ValueOf(pattern).Kind()

			if valueType == reflect.Map && patternType == reflect.Map {
				log.Debugf("recursing into map %v", currentFieldPath)

				result = jc.doCompareJSONRegexUnmarshaled(
					pattern.(map[string]interface{}),
					value.(map[string]interface{}),
					optionalPaths,
					currentFieldPath,
				)

				if !result {
					return false
				}
			} else if (valueType == reflect.Array || valueType == reflect.Slice) &&
				(patternType == reflect.Array || patternType == reflect.Slice) {

				log.Debugf("recursing into array %v", currentFieldPath)

				valueJsonBytes, err1 := json.Marshal(value)
				patternJsonBytes, err2 := json.Marshal(pattern)

				if err1 != nil || err2 != nil {
					log.Errorf("value %v raised %v and pattern %v raised %v",
						value, err1, pattern, err2)
					return false
				}

				result = jc.doCompareArrayRegex(
					string(patternJsonBytes), string(valueJsonBytes), optionalPaths, currentFieldPath)

				if !result {
					return false
				}
			} else if eql := reflect.DeepEqual(pattern, value); !eql {
				log.Debugf("value %v doesn't DeepEqual %v", value, pattern)
				return false
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

func (jc *JSONComparator) Compare(
	s1, s2 string,
	optionalPaths map[string]bool,
	currentPath string) bool {

	if isArray(s1) != isArray(s2) {
		log.Debugf("only one of these is an array %v %v", s1, s2)
		return false
	}

	if isArray(s1) || isArray(s2) {
		return jc.doCompareArrayRegex(s1, s2, optionalPaths, currentPath)
	}

	return jc.doCompareJSONRegex(s1, s2, optionalPaths, currentPath)
}
