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

func (pv PersistVars) Fill(m *definition.Mock, input string) string {
	r := regexp.MustCompile(`\{\{\s*persist\.(.+?)\s*\}\}`)

	return r.ReplaceAllStringFunc(input, func(raw string) string {
		// replace the strings
		if r, found := pv.replaceString(m, raw); found {
			return r
		}
		// replace regexes
		return raw
	})

}

func (pv PersistVars) replaceString(m *definition.Mock, raw string) (string, bool) {
	found := false
	s := ""
	tag := strings.Trim(raw[2:len(raw)-2], " ")
	if tag == "persist.entity.name" {
		s = m.Persist.Entity
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
	}

	if !found {
		return raw, false
	}
	return s, true
}
