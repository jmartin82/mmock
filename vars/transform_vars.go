package vars

import (
	"regexp"
	"strings"

	"github.com/jmartin82/mmock/definition"
)

type TransformVars struct {
}

func (tv TransformVars) Fill(m *definition.Mock, input string) string {
	r := regexp.MustCompile(`\{\{\s*transform\.(.+?)\s*\}\}`)

	return r.ReplaceAllStringFunc(input, func(raw string) string {
		// replace the strings
		if r, found := tv.replaceString(m, raw); found {
			return r
		}
		// replace regexes
		return raw
	})

}

func (tv TransformVars) replaceString(m *definition.Mock, raw string) (string, bool) {
	found := false
	s := ""
	tag := strings.Trim(raw[2:len(raw)-2], " ")
	if i := strings.Index(tag, "transform.merge"); i == 0 {
		s, found = tv.mergeContent(m, tag[15:])
	}

	if !found {
		return raw, false
	}
	return s, true
}

func (tv TransformVars) mergeContent(m *definition.Mock, name string) (string, bool) {
	r := regexp.MustCompile(`\((.*)\)`)
	if len(r.FindStringSubmatch(name)) > 0 {
		parts := strings.Split(r.FindStringSubmatch(name)[1], ",")
		return strings.Join(parts, "#"), true
	}
	return "", false
}
