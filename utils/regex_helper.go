package utils

import "regexp"

type RegexHelper struct {
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
