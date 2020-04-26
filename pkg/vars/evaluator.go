package vars

import (
	"regexp"
	"strings"

	"github.com/jmartin82/mmock/pkg/mock"
)

var varsRegex = regexp.MustCompile(`\{\{\s*(.+?)\s*\}\}`)

type Evaluator interface {
	Eval(req *mock.Request, m *mock.Definition)
}

type ResponseMessageEvaluator struct {
	FillerFactory FillerFactory
}

func NewResponseMessageEvaluator(fp FillerFactory) *ResponseMessageEvaluator {
	return &ResponseMessageEvaluator{FillerFactory: fp}
}

func (fp ResponseMessageEvaluator) Eval(req *mock.Request, m *mock.Definition) {
	requestFiller := fp.FillerFactory.CreateRequestFiller(req, m)
	fakeFiller := fp.FillerFactory.CreateFakeFiller()
	streamFiller := fp.FillerFactory.CreateStreamFiller()
	holders := fp.walkAndGet(m.Response)
	fp.extractVars(m.Callback.Body, &holders)

	vars := requestFiller.Fill(holders)
	fp.mergeVars(vars, fakeFiller.Fill(holders))
	fp.mergeVars(vars, streamFiller.Fill(holders))
	fp.walkAndFill(m, vars)
}

func (fp ResponseMessageEvaluator) walkAndGet(res mock.Response) []string {

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

func (fp ResponseMessageEvaluator) walkAndFill(m *mock.Definition, vars map[string][]string) {
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
	m.Callback.Body = fp.replaceVars(m.Callback.Body, vars)
}

func (fp ResponseMessageEvaluator) replaceVars(input string, vars map[string][]string) string {
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

func (fp ResponseMessageEvaluator) extractVars(input string, vars *[]string) {
	if m := varsRegex.FindAllString(input, -1); m != nil {
		for _, v := range m {
			varName := strings.Trim(v, "{} ")
			*vars = append(*vars, varName)
		}
	}
}

func (fp ResponseMessageEvaluator) mergeVars(org map[string][]string, vals map[string][]string) {
	for k, v := range vals {
		org[k] = v
	}
}
