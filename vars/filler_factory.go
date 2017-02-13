package vars

import (
	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/vars/fakedata"
)

type FillerFactory interface {
	CreateRequestFiller(req *definition.Request, mock *definition.Mock) Filler
	CreateFakeFiller() Filler
}

type MockFillerFactory struct {
	FakeAdapter fakedata.DataFaker
}

func (mff MockFillerFactory) CreateRequestFiller(req *definition.Request, mock *definition.Mock) Filler {
	return Request{Mock: mock, Request: req}
}

func (mff MockFillerFactory) CreateFakeFiller() Filler {

	return Fake{Fake: mff.FakeAdapter}
}
