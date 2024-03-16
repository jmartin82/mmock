package payload

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"
)

var DEBUG bool

type JSONComparator struct {
}

func isArray(s string) bool {
	st := strings.TrimLeft(s, " ")
	return len(st) > 0 && st[0] == '['
}

func (jc *JSONComparator) doCompareJSONRegexUnmarshaled(patterns, values map[string]interface{}) bool {
	var matches bool
	matches = jc.match(patterns, values)
	if !matches && DEBUG {
		log.Printf("values: %v don't match: %v", values, patterns)
	}
	return matches
}

func (jc *JSONComparator) doCompareJSONRegex(jsonWithPatterns, jsonWithValues string) bool {
	var patterns map[string]interface{}
	var values map[string]interface{}
	if err := json.Unmarshal([]byte(jsonWithPatterns), &patterns); err != nil {
		if DEBUG {
			log.Printf("error in json patterns: %v", err)
		}

		return false
	}
	if err := json.Unmarshal([]byte(jsonWithValues), &values); err != nil {
		if DEBUG {
			log.Printf("error in json values: %v", err)
		}

		return false
	}
	return jc.doCompareJSONRegexUnmarshaled(patterns, values)
}

func (jc *JSONComparator) doCompareArrayRegex(jsonWithPatterns, jsonWithValues string) bool {
	var patterns []map[string]interface{}
	var values []map[string]interface{}

	if err := json.Unmarshal([]byte(jsonWithPatterns), &patterns); err != nil {
		if DEBUG {
			log.Printf("error in json patterns: %v", err)
		}

		return false
	}
	if err := json.Unmarshal([]byte(jsonWithValues), &values); err != nil {
		if DEBUG {
			log.Printf("error in json patterns: %v", err)
		}

		return false
	}
	return jc.doCompareArrayRegexUnmarshaled(patterns, values)
}

func (jc *JSONComparator) doCompareArrayRegexUnmarshaled(patterns, values []map[string]interface{}) bool {

	for i := 0; i < len(patterns); i++ {
		if !jc.match(patterns[i], values[i]) {
			if DEBUG {
				log.Printf("value %v doesn't match %v", values[i], patterns[i])
			}

			return false
		}
	}
	return true
}

func (jc *JSONComparator) match(p, v map[string]interface{}) bool {
	for field, pattern := range p {

		value, exists := v[field]
		if DEBUG {
			log.Printf("comparing field %v with pattern %v against value %v", field, pattern, value)
		}

		if !exists {
			if DEBUG {
				log.Printf("field doesn't exist: %v", field)
			}

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
				if DEBUG {
					log.Printf("recursing into map %v", field)
				}

				result = jc.doCompareJSONRegexUnmarshaled(
					pattern.(map[string]interface{}),
					value.(map[string]interface{}))

				if !result {
					return false
				}
			} else if (valueType == reflect.Array || valueType == reflect.Slice) &&
				(patternType == reflect.Array || patternType == reflect.Slice) {

				if DEBUG {
					log.Printf("recursing into array %v", field)
				}

				result = jc.doCompareArrayRegexUnmarshaled(
					pattern.([]map[string]interface{}),
					value.([]map[string]interface{}))

				if !result {
					return false
				}
			} else {
				var eql bool
				eql = reflect.DeepEqual(pattern, value)

				if !eql {
					if DEBUG {
						log.Printf("value %v doesn't DeepEqual %v", value, pattern)
					}
				}

				return eql
			}
		}
		matched, err := regexp.MatchString(str, fmt.Sprint(value))
		if err != nil || !matched {
			if DEBUG {
				log.Printf("value %v doesn't match %v : %v", fmt.Sprint(value), str, err)
			}

			return false
		}
	}
	return true
}

func (jc *JSONComparator) Compare(s1, s2 string) bool {
	DEBUG = os.Getenv("DEBUG") == "true"

	if isArray(s1) != isArray(s2) {
		if DEBUG {
			log.Printf("only one of these is an array %v %v", s1, s2)
		}

		return false
	}

	if isArray(s1) || isArray(s2) {
		return jc.doCompareArrayRegex(s1, s2)
	}

	return jc.doCompareJSONRegex(s1, s2)
}
