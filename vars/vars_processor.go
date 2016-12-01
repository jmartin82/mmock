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
	requestFiller := fp.FillerFactory.CreateRequestFiller(req)
	fakeFiller := fp.FillerFactory.CreateFakeFiller(fp.FakeAdapter)
	storageFiller := fp.FillerFactory.CreateStorageFiller(fp.PersistEngines)
	persistFiller := fp.FillerFactory.CreatePersistFiller(fp.PersistEngines)
	entityActions := persist.EntityActions{fp.PersistEngines}

	fp.walkAndFill(requestFiller, m, true)
	fp.walkAndFill(fakeFiller, m, true)
	fp.walkAndFill(storageFiller, m, true)

	// we need to make sure the persisted vars are filled before executing the actions - as we need to make sure the persist vars are replaced in the persist actions
	fp.walkAndFillPersisted(persistFiller, m)

	entityActions.ApplyActions(m)

	fp.walkAndFill(persistFiller, m, false)

}

func (fp VarsProcessor) walkAndFill(f Filler, m *definition.Mock, fillPersisted bool) {
	res := &m.Response
	amqp := &m.Notify.Amqp
	for header, values := range res.Headers {
		for i, value := range values {
			res.Headers[header][i] = f.Fill(m, value, false)
		}

	}
	for cookie, value := range res.Cookies {
		res.Cookies[cookie] = f.Fill(m, value, false)
	}

	amqp.Body = f.Fill(m, amqp.Body, false)
	res.Body = f.Fill(m, res.Body, false)

	if fillPersisted {
		fp.walkAndFillPersisted(f, m)
	}
}

func (fp VarsProcessor) walkAndFillPersisted(f Filler, m *definition.Mock) {
	per := &m.Persist

	per.Entity = f.Fill(m, per.Entity, false)
	per.Collection = f.Fill(m, per.Collection, true)
	for action, value := range per.Actions {
		per.Actions[action] = f.Fill(m, value, false)
	}
}
