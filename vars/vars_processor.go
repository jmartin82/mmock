package vars

import (
	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/persist"
	"github.com/jmartin82/mmock/vars/fakedata"
)

type VarsProcessor struct {
	FillerFactory  FillerFactory
	FakeAdapter    fakedata.DataFaker
	PersistEngines *persist.PersistEngineBag
}

func (fp VarsProcessor) Eval(req *definition.Request, m *definition.Mock) {
	fp.walkAndFill(fp.FillerFactory.CreateRequestFiller(req), m)
	fp.walkAndFill(fp.FillerFactory.CreateFakeFiller(fp.FakeAdapter), m)
	fp.walkAndFill(fp.FillerFactory.CreatePersistFiller(fp.PersistEngines), m)
}

func (fp VarsProcessor) walkAndFill(f Filler, m *definition.Mock) {
	res := &m.Response
	per := &m.Persist
	amqp := &m.Notify.Amqp
	for header, values := range res.Headers {
		for i, value := range values {
			res.Headers[header][i] = f.Fill(m, value)
		}

	}
	for cookie, value := range res.Cookies {
		res.Cookies[cookie] = f.Fill(m, value)
	}

	res.Body = f.Fill(m, res.Body)
	amqp.Body = f.Fill(m, amqp.Body)
	per.Entity = f.Fill(m, per.Entity)
	for action, value := range per.Actions {
		per.Actions[action] = f.Fill(m, value)
	}

}
