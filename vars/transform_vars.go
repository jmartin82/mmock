package vars

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/Jeffail/gabs"
)

type TransformVars struct {
}

func (tv TransformVars) Operate(holders []string, vars map[string]string) map[string]string {

	result := make(map[string]string)

	for _, tag := range holders {
		found := false
		s := ""
		if i := strings.Index(tag, "transform.join"); i == 0 {
			s, found = tv.join(tag[14:], vars)
		}
		if found {
			result[tag] = s
		}
	}
	return result
}

func (tv TransformVars) GetFunctionVars(vars []string) []string {
	var result = []string{}
	for _, name := range vars {
		parms := tv.getParams(name)
		result = append(result, parms...)
	}
	return result
}

func (tv TransformVars) getParams(name string) []string {
	r := regexp.MustCompile(`\((.*)\)`)
	if m := r.FindStringSubmatch(name); len(m) > 0 {
		return strings.Split(m[1], ",")
	}
	return []string{}
}

func (tv TransformVars) getReplacedParams(name string, value map[string]string) []string {
	parms := tv.getParams(name)
	replaced := []string{}

	for _, param := range parms {
		if v, f := value[param]; f {
			replaced = append(replaced, v)
		} else {
			replaced = append(replaced, strings.Trim(param, "\"'"))
		}

	}
	return replaced

}

func (tv TransformVars) joinText(name string, value map[string]string) (string, bool) {
	inputs := tv.getReplacedParams(name, value)
	if len(inputs) < 2 {
		return "", false
	}
	return strings.Join(inputs, ""), true
}

func (tv TransformVars) isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func (tv TransformVars) join(name string, value map[string]string) (string, bool) {
	inputs := tv.getReplacedParams(name, value)
	if len(inputs) < 1 {
		return "", false
	}
	if len(inputs) < 2 {
		return inputs[0], true
	}
	result := gabs.New()
	text := ""
	for _, input := range inputs {
		if tv.isJSON(input) {
			jsonParsed, _ := gabs.ParseJSON([]byte(input))
			children, _ := jsonParsed.S().ChildrenMap()

			for key, child := range children {
				result.Set(child.Data(), key)
			}
		} else {
			text += input
			if tv.isJSON(text) {
				jsonParsed, _ := gabs.ParseJSON([]byte(text))
				children, _ := jsonParsed.S().ChildrenMap()

				for key, child := range children {
					result.Set(child.Data(), key)
				}
				text = ""
			}
		}

	}

	if text != "" {
		return strings.Join(inputs, ""), true
	}

	return result.String(), true
}
