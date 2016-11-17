package vars

import (
	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/persist"
	"github.com/jmartin82/mmock/vars/fakedata"
)

type FillerFactory interface {
	CreateRequestFiller(req *definition.Request) Filler
	CreateFakeFiller(Fake fakedata.DataFaker) Filler
	CreatePersistFiller(Engines *persist.PersistEngineBag) Filler
}

type MockFillerFactory struct{}

func (mff MockFillerFactory) CreateRequestFiller(req *definition.Request) Filler {
	return RequestVars{Request: req}
}
func (mff MockFillerFactory) CreateFakeFiller(fake fakedata.DataFaker) Filler {
	return FakeVars{Fake: fake}
}
func (mff MockFillerFactory) CreatePersistFiller(engines *persist.PersistEngineBag) Filler {
	return PersistVars{Engines: engines}
}
