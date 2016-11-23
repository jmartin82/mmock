package vars

import (
	"regexp"
	"strings"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/persist"
)

type PersistVars struct {
	Engines *persist.PersistEngineBag
}

func (pp PersistVars) Fill(m *definition.Mock, input string) string {
	r := regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_\.]+)\s*\}\}`)
	return r.ReplaceAllStringFunc(input, func(raw string) string {
		found := false
		s := ""
		tag := strings.Trim(raw[2:len(raw)-2], " ")
		if i := strings.Index(tag, "persist.entity.name"); i == 0 {
			s = m.Persist.Entity
			found = true
		} else if i := strings.Index(tag, "persist.entity.content"); i == 0 {
			engine := pp.Engines.Get(m.Persist.Engine)
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
			return raw
		}
		return s
	})
}
