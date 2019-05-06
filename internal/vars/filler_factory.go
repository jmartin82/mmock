package vars

import (
	"github.com/jmartin82/mmock/pkg/mock"
	"github.com/jmartin82/mmock/internal/vars/fakedata"
)

type FillerFactory interface {
	CreateRequestFiller(req *mock.Request, mock *mock.Definition) Filler
	CreateFakeFiller() Filler
	CreateStreamFiller() Filler
}

type MockFillerFactory struct {
	FakeAdapter fakedata.DataFaker
}

func (mff MockFillerFactory) CreateRequestFiller(req *mock.Request, mock *mock.Definition) Filler {
	return Request{Mock: mock, Request: req}
}

func (mff MockFillerFactory) CreateFakeFiller() Filler {

	return Fake{Fake: mff.FakeAdapter}
}

func (mff MockFillerFactory) CreateStreamFiller() Filler {
	return Stream{}
}
