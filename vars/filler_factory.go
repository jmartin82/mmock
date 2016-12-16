package vars

import (
	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/vars/fakedata"
)

type FillerFactory interface {
	CreateRequestFiller(req *definition.Request, mock *definition.Mock) Filler
	CreateFakeFiller(Fake fakedata.DataFaker) Filler
}

type MockFillerFactory struct{}

func (mff MockFillerFactory) CreateRequestFiller(req *definition.Request, mock *definition.Mock) Filler {
	return RequestVars{Mock: mock, Request: req}
}

func (mff MockFillerFactory) CreateFakeFiller(fake fakedata.DataFaker) Filler {
	return FakeVars{Fake: fake}
}
