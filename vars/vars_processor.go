package vars

import (
	"regexp"

	"log"

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
	persistFiller := fp.FillerFactory.CreatePersistFiller(fp.PersistEngines)
	transformFiller := fp.FillerFactory.CreateTransformFiller()
	entityActions := persist.EntityActions{fp.PersistEngines}

	log.Fatalln(fp.walkAndGet(m.Response))

	fp.walkAndFill(requestFiller, m)
	fp.walkAndFill(fakeFiller, m)
	entityActions.ApplyActions(m)
	fp.walkAndFill(persistFiller, m)
	fp.walkAndFill(transformFiller, m)

}

func (fp VarsProcessor) walkAndGet(res definition.Response) []string {
	r := regexp.MustCompile(`\{\{\s*(.+?)\s*\}\}`)
	vars := []string{}
	for _, header := range res.Headers {
		for _, value := range header {
			m := r.FindStringSubmatch(value)
			if len(m) > 1 {
				vars = append(vars, m[1:]...)
			}
		}

	}
	for _, value := range res.Cookies {

		m := r.FindStringSubmatch(value)
		if len(m) > 1 {
			vars = append(vars, m[1:]...)
		}
	}

	if m := r.FindAllString(res.Body, -1); m != nil {
		vars = append(vars, m[1])
	}

	return vars
}

func (fp VarsProcessor) walkAndFill(f Filler, m *definition.Mock) {
	res := &m.Response
	per := &m.Persist
	for header, values := range res.Headers {
		for i, value := range values {
			res.Headers[header][i] = f.Fill(m, value)
		}

	}
	for cookie, value := range res.Cookies {
		res.Cookies[cookie] = f.Fill(m, value)
	}

	res.Body = f.Fill(m, res.Body)
	per.Entity = f.Fill(m, per.Entity)
	for action, value := range per.Actions {
		per.Actions[action] = f.Fill(m, value)
	}

}
