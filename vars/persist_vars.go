package vars

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/persist"
)

type PersistVars struct {
	Engines   *persist.PersistEngineBag
	processed bool
}

func (pp PersistVars) getEngine(m *definition.Mock) persist.EntityPersister {
	engine := pp.Engines.Get(m.Persist.Engine)
	//fix persiste console logging
	if m.Persist.Engine == "" {
		m.Persist.Engine = engine.GetName()
	}
	return engine
}

func (pp PersistVars) Fill(m *definition.Mock, input string) string {
	if !pp.processed {
		pp.ApplyActions(m)
	}
	r := regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_\.]+)\s*\}\}`)
	engine := pp.getEngine(m)

	return r.ReplaceAllStringFunc(input, func(raw string) string {
		found := false
		s := ""
		tag := strings.Trim(raw[2:len(raw)-2], " ")
		if i := strings.Index(tag, "persist.entity.name"); i == 0 {
			s = m.Persist.Entity
			found = true
		} else if i := strings.Index(tag, "persist.entity.content"); i == 0 {
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

func (pp PersistVars) ApplyActions(m *definition.Mock) {
	engine := pp.getEngine(m)
	fileName := m.Persist.Entity
	if value, ok := m.Persist.Actions["write"]; ok {
		if err := engine.Write(fileName, value); err != nil {
			log.Println("Error writing in a entity")
			return
		}
	}

	if value, ok := m.Persist.Actions["append"]; ok {
		content, err := engine.Read(fileName)
		if err != nil {
			log.Println("Error reading in a entity")
			return
		}
		if isJSON(content) && isJSON(value) {
			content = joinJSON(content, value)
		} else if isJSON(content) && !isJSON(value) {
			log.Printf("There is no way to append this : %s\n", value)
		} else {
			content += value
		}
		if err := engine.Write(fileName, content); err != nil {
			log.Println("Error appending in a entity")
		}
	}

	if _, ok := m.Persist.Actions["delete"]; ok {
		if err := engine.Delete(fileName); err != nil {
			log.Println("Error deleting a entity")
			return
		}
	}
}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func joinJSON(inputs ...string) string {
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
