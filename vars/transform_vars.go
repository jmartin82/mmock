package vars

import (
	"bytes"
	"regexp"
	"strings"
)

type TransformVars struct {
}

func (tv TransformVars) Operate(holders []string, vars map[string]string) map[string]string {

	result := make(map[string]string)

	for _, tag := range holders {
		found := false
		s := ""
		if i := strings.Index(tag, "transform.merge"); i == 0 {
			s, found = tv.mergeContent(tag[15:], vars)
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

func (tv TransformVars) mergeContent(name string, value map[string]string) (string, bool) {
	found := false
	parms := tv.getParams(name)
	var buffer bytes.Buffer

	for _, param := range parms {
		if v, f := value[param]; f {
			buffer.WriteString(v)
			found = true
		} else {
			//regular string
			buffer.WriteString(strings.Trim(param, "\"'"))
		}

	}
	return buffer.String(), found

}
