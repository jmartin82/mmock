package vars

import (
	"regexp"
	"strings"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/persist"
	"github.com/jmartin82/mmock/utils"
)

type PersistVars struct {
	Engines     *persist.PersistEngineBag
	RegexHelper utils.RegexHelper
}

func (pv PersistVars) Fill(m *definition.Mock, input string, multipleMatch bool) string {
	r := regexp.MustCompile(`\{\{(.+?)\}\}`)

	if !multipleMatch {
		return r.ReplaceAllStringFunc(input, func(raw string) string {
			// replace the strings
			if raw, found := pv.replaceString(m, raw); found {
				return raw
			}
			// replace regexes
			return pv.replaceRegex(m, raw)
		})
	} else {
		// first replace all strings
		input = r.ReplaceAllStringFunc(input, func(raw string) string {
			item, _ := pv.replaceString(m, raw)
			return item
		})
		// get multiple entities using regex
		results, found := pv.RegexHelper.GetCollectionItems(input, func(raw string) (string, string, bool) {
			return pv.getPersistRegexParts(m, raw)
		})
		if found {
			if len(results) == 1 {
				return "," + results[0] // add a comma in the beginning so that we will now that the item is a single entity
			}

			return strings.Join(results, ",")
		}
		return input
	}
}

func (pv PersistVars) replaceString(m *definition.Mock, raw string) (string, bool) {
	found := false
	s := ""
	tag := strings.Trim(raw[2:len(raw)-2], " ")
	if tag == "persist.entity.name" {
		s = m.Persist.Entity
		found = true
	} else if i := strings.Index(tag, "persist.collection.name"); i == 0 {
		s = m.Persist.Collection
		found = true
	} else if i := strings.Index(tag, "persist.entity.content"); i == 0 {
		engine := pv.Engines.Get(m.Persist.Engine)
		content, err := engine.Read(m.Persist.Entity)
		//if error, we change Response status and body
		if err != nil {
			s = ""
			m.Response.Body = ""
			m.Response.StatusCode = 404
		}
		s = content
		found = true
	} else if i := strings.Index(tag, "persist.collection.content"); i == 0 {
		engine := pv.Engines.Get(m.Persist.Engine)
		content, err := engine.ReadCollection(m.Persist.Collection)
		//if error, we change Response status and body
		if err != nil {
			s = ""
			m.Response.Body = ""
			m.Response.StatusCode = 404
		}
		s = content
		found = true
	}

	if !found {
		return raw, false
	}
	return s, true
}

func (rp PersistVars) getPersistRegexParts(m *definition.Mock, input string) (string, string, bool) {
	if i := strings.Index(input, "persist.entity.name"); i == 0 && len(input) > len("persist.entity.name") {
		return m.Persist.Entity, input[20:], true
	}
	return "", "", false
}

func (rp PersistVars) replaceRegex(m *definition.Mock, raw string) string {
	tag := strings.Trim(raw[2:len(raw)-2], " ")
	if regexInput, regexPattern, found := rp.getPersistRegexParts(m, tag); found {
		if result, found := rp.RegexHelper.GetStringPart(regexInput, regexPattern, "value"); found {
			return result
		}
	}
	return raw
}
