package vars

import (
	"github.com/jmartin82/mmock/pkg/mock"
	"regexp"
	"strings"


)

var varsRegex = regexp.MustCompile(`\{\{\s*(.+?)\s*\}\}`)

type Evaluator struct {
	FillerFactory FillerFactory
}

func (fp Evaluator) Eval(req *mock.Request, m *mock.Definition) {
	requestFiller := fp.FillerFactory.CreateRequestFiller(req, m)
	fakeFiller := fp.FillerFactory.CreateFakeFiller()
	streamFiller := fp.FillerFactory.CreateStreamFiller()
	holders := fp.walkAndGet(m.Response)

	vars := requestFiller.Fill(holders)
	fp.mergeVars(vars, fakeFiller.Fill(holders))
	fp.mergeVars(vars, streamFiller.Fill(holders))
	fp.walkAndFill(m, vars)
}

func (fp Evaluator) walkAndGet(res mock.Response) []string {

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

func (fp Evaluator) walkAndFill(m *mock.Definition, vars map[string][]string) {
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

func (fp Evaluator) replaceVars(input string, vars map[string][]string) string {
	return varsRegex.ReplaceAllStringFunc(input, func(value string) string {
		varName := strings.Trim(value, "{} ")
		// replace the strings
		if v, found := vars[varName]; found {
			r := v[0]
			vars[varName] = v[1:]
			return r
		}
		// replace regexes
		return value
	})
}

func (fp Evaluator) extractVars(input string, vars *[]string) {
	if m := varsRegex.FindAllString(input, -1); m != nil {
		for _, v := range m {
			varName := strings.Trim(v, "{} ")
			*vars = append(*vars, varName)
		}
	}
}

func (fp Evaluator) mergeVars(org map[string][]string, vals map[string][]string) {
	for k, v := range vals {
		org[k] = v
	}
}
