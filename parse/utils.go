package parse

import (
	"encoding/json"
	"regexp"

	"strings"

	"github.com/Jeffail/gabs"
	"github.com/jmartin82/mmock/definition"
)

//IsJSON checks if a string is valid JSON or not
func IsJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

//JoinJSON merges two JSON strings
func JoinJSON(inputs ...string) string {
	if len(inputs) == 1 {
		return inputs[0]
	}

	result := gabs.New()
	for _, input := range inputs {
		jsonParsed, _ := gabs.ParseJSON([]byte(input))
		children, _ := jsonParsed.S().ChildrenMap()

		for key, child := range children {
			result.Set(child.Data(), key)
		}
	}

	return result.String()
}

//GetStringPart gets regex mached group value for a given input and regex pattern
func GetStringPart(input string, pattern string, groupName string) (string, bool) {
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

//FormatJSON formats a JSON string
func FormatJSON(input string) (result string, err error) {
	var jsonParsed interface{}
	json.Unmarshal([]byte(input), &jsonParsed)
	if err != nil {
		return "", err
	}

	byteString, err := json.Marshal(jsonParsed)
	if err != nil {
		return "", err
	}

	return string(byteString), nil
}

//JSONSStringsAreEqual checks whether two JSON strings are actually equal JSON objects
func JSONSStringsAreEqual(input1 string, input2 string) (result bool, err error) {
	formatedInput1, err := FormatJSON(input1)
	if err != nil {
		return false, err
	}
	formatedInput2, err := FormatJSON(input2)
	if err != nil {
		return false, err
	}
	return formatedInput1 == formatedInput2, nil
}

//WrapNonJSONStringIfNeeded wrapps non JSON string in NonJSONItem object
func WrapNonJSONStringIfNeeded(input string) (result string, err error) {
	if !IsJSON(input) {
		wrapper := definition.NonJSONItem{Content: input}
		bytesString, err := json.Marshal(wrapper)
		if err != nil {
			return "", err
		}
		return string(bytesString), nil
	}
	return input, nil
}

//UnWrapNonJSONStringIfNeeded wrapps non JSON string in NonJSONItem object
func UnWrapNonJSONStringIfNeeded(input string) string {
	if IsJSON(input) && strings.Index(input, "non_json_content") > -1 {
		var nonJSON definition.NonJSONItem
		err := json.Unmarshal([]byte(input), &nonJSON)
		if err != nil || nonJSON.Content == "" { // the json most probably is not a NonJSONItem
			return input
		}

		return nonJSON.Content
	}
	return input
}
