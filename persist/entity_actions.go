package persist

import (
	"log"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/utils"
)

//FilePersister persists body in file
type EntityActions struct {
	Engines *PersistEngineBag
}

func (ea EntityActions) getEngine(m *definition.Mock) EntityPersister {
	engine := ea.Engines.Get(m.Persist.Engine)
	//fix persiste console logging
	if m.Persist.Engine == "" {
		m.Persist.Engine = engine.GetName()
	}
	return engine
}

func (ea EntityActions) ApplyActions(m *definition.Mock) {
	engine := ea.getEngine(m)
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
		if utils.IsJSON(content) && utils.IsJSON(value) {
			content = utils.JoinJSON(content, value)
		} else if utils.IsJSON(content) && !utils.IsJSON(value) {
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
