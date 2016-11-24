package utils

import (
	"regexp"
	"strings"
)

type RegexHelper struct {
}

func (rh RegexHelper) GetCollectionItems(input string, getRegexParts func(string) (string, string, bool)) ([]string, bool) {
	matchMap := rh.getMatchMap(input, getRegexParts)
	if len(matchMap) == 0 {
		return []string{}, false
	}

	cartesian := Cartesian{}
	combinationItems := cartesian.GetCombinations(matchMap)

	results := []string{}

	for _, currentMap := range combinationItems {
		resultItem := input
		for key, value := range currentMap {
			resultItem = strings.Replace(resultItem, key, value, -1)
		}
		results = append(results, resultItem)
	}
	return results, true
}

func (rh RegexHelper) getMatchMap(input string, getRegexParts func(string) (string, string, bool)) map[string][]string {
	regexMap := make(map[string][]string)
	r := regexp.MustCompile(`\{\{(.+?)\}\}`)

	for _, regexPattern := range r.FindAllString(input, -1) {
		if _, exists := regexMap[regexPattern]; !exists {
			regexMap[regexPattern] = []string{regexPattern}
		}
	}

	for key, _ := range regexMap {
		pattern := strings.Trim(key[2:len(key)-2], " ")
		if regexInput, regexPattern, found := getRegexParts(pattern); found {
			if calculatedValues, found := rh.getStringParts(regexInput, regexPattern, "value"); found {
				regexMap[key] = calculatedValues
			} else {
				regexMap[key] = []string{key}
			}
		} else {
			regexMap[key] = []string{key}
		}
	}
	return regexMap
}

func (rh RegexHelper) getStringParts(input string, pattern string, groupName string) ([]string, bool) {
	r, error := regexp.Compile(pattern)
	if error != nil {
		return []string{}, false
	}

	names := []string{}

	matches := r.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		value, found := rh.getValue(r, match, "value")
		if found {
			names = append(names, value)
		}
	}
	return names, len(names) > 0
}

func (rh RegexHelper) getValue(r *regexp.Regexp, match []string, groupName string) (value string, present bool) {
	result := make(map[string]string)
	names := r.SubexpNames()
	if len(match) >= len(names) {
		for i, name := range names {
			if i != 0 {
				result[name] = match[i]
			}
		}
	}

	value, present = result[groupName]

	return value, present
}

// GetStringPart gets the value of the group name matching the input using the pattern
func (rh RegexHelper) GetStringPart(input string, pattern string, groupName string) (string, bool) {
	r, error := regexp.Compile(pattern)
	if error != nil {
		return "", false
	}

	match := r.FindStringSubmatch(input)
	result := make(map[string]string)
	names := r.SubexpNames()
	if len(match) >= len(names) {
		for i, name := range names {
			if i != 0 {
				result[name] = match[i]
			}
		}
	}

	value, present := result[groupName]

	return value, present
}
