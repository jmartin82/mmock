package vars

import (
	"regexp"
	"strings"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/persist"
	"github.com/jmartin82/mmock/vars/fakedata"
)

var varsRegex = regexp.MustCompile(`\{\{\s*(.+?)\s*\}\}`)

type VarsProcessor struct {
	FillerFactory  FillerFactory
	FakeAdapter    fakedata.DataFaker
	PersistEngines *persist.PersistEngineBag
}

func (fp VarsProcessor) Eval(req *definition.Request, m *definition.Mock) {
	requestFiller := fp.FillerFactory.CreateRequestFiller(req, m)
	fakeFiller := fp.FillerFactory.CreateFakeFiller(fp.FakeAdapter)
	transformFiller := fp.FillerFactory.CreateTransformFiller()

	//get var holders
	holders := fp.walkAndGet(m.Response)
	inputParams := transformFiller.GetFunctionVars(holders)
	holders = append(holders, inputParams...)

	//fill vars
	vars := requestFiller.Fill(holders)
	fp.mergeVars(vars, fakeFiller.Fill(holders))
	fp.mergeVars(vars, transformFiller.Operate(holders, vars))

	//replace vars
	fp.walkAndFill(m, vars)
}

func (fp VarsProcessor) walkAndGet(res definition.Response) []string {

	vars := []string{}
	for _, header := range res.Headers {
		for _, value := range header {
			fp.extractVars(value, &vars)
		}

	}
	for _, value := range res.Cookies {
		fp.extractVars(value, &vars)
	}

	fp.extractVars(res.Body, &vars)
	return vars
}

func (fp VarsProcessor) walkAndFill(m *definition.Mock, vars map[string]string) {
	res := &m.Response
	for header, values := range res.Headers {
		for i, value := range values {
			res.Headers[header][i] = fp.replaceVars(value, vars)
		}

	}
	for cookie, value := range res.Cookies {
		res.Cookies[cookie] = fp.replaceVars(value, vars)
	}

	res.Body = fp.replaceVars(res.Body, vars)
}

func (fp VarsProcessor) replaceVars(input string, vars map[string]string) string {
	return varsRegex.ReplaceAllStringFunc(input, func(value string) string {
		varName := strings.Trim(value, "{} ")
		// replace the strings
		if r, found := vars[varName]; found {
			return r
		}
		// replace regexes
		return value
	})
}

func (fp VarsProcessor) extractVars(input string, vars *[]string) {
	if m := varsRegex.FindAllString(input, -1); m != nil {
		for _, v := range m {
			varName := strings.Trim(v, "{} ")
			*vars = append(*vars, varName)
		}
	}
}

func (fp VarsProcessor) mergeVars(org map[string]string, vals map[string]string) {
	for k, v := range vals {
		org[k] = v
	}
}
