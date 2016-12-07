package persist

import (
	"strings"

	"github.com/jmartin82/mmock/logging"
)

type PersistEngineBag struct {
	engines map[string]EntityPersister
}

func (peb *PersistEngineBag) SetDefault(def EntityPersister) {
	peb.engines["default"] = def
	peb.Add(def)
}

func (peb *PersistEngineBag) Add(engine EntityPersister) {
	name := strings.ToLower(engine.GetName())
	peb.engines[name] = engine
}

func (peb *PersistEngineBag) Get(name string) EntityPersister {
	name = strings.ToLower(name)
	if engine, ok := peb.engines[name]; ok {
		return engine
	}

	def := peb.engines["default"]
	logging.Printf("Using the default persist engine: %s\n", def.GetName())
	return def
}

func GetNewPersistEngineBag(def EntityPersister) *PersistEngineBag {
	bag := make(map[string]EntityPersister)
	p := PersistEngineBag{engines: bag}
	p.SetDefault(def)
	return &p
}
