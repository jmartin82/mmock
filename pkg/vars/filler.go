package vars

import (
	"github.com/jmartin82/mmock/pkg/vars/fake"
	"github.com/jmartin82/mmock/pkg/mock"
)

type Filler interface {
	Fill(holders []string) map[string][]string
}

type FillerFactory interface {
	CreateRequestFiller(req *mock.Request, mock *mock.Definition) Filler
	CreateFakeFiller() Filler
	CreateStreamFiller() Filler
}

type MockFillerFactory struct {
	FakeDataProvider fake.Generator
}

func NewFillerFactory(fdp fake.Generator) *MockFillerFactory {
	return &MockFillerFactory{FakeDataProvider: fdp}
}

func (mff MockFillerFactory) CreateRequestFiller(req *mock.Request, mock *mock.Definition) Filler {
	return Request{Mock: mock, Request: req}
}

func (mff MockFillerFactory) CreateFakeFiller() Filler {

	return Fake{Fake: mff.FakeDataProvider}
}

func (mff MockFillerFactory) CreateStreamFiller() Filler {
	return Stream{}
}
