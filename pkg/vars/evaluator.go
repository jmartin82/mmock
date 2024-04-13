package vars

import (
	"regexp"
	"strings"

	"github.com/jmartin82/mmock/v3/pkg/mock"
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

	//replace stream holders for their content
	m.Response.HttpHeaders, m.Response.Body = fp.walkAndFill(m.Response, streamFiller.Fill(fp.walkAndGet(m.Response)))
	m.Callback.HttpHeaders, m.Callback.Body = fp.walkAndFill(m.Callback, streamFiller.Fill(fp.walkAndGet(m.Callback)))

	//extract holders
	holders := fp.walkAndGet(m.Response)
	holders = append(holders, fp.walkAndGet(m.Callback)...)

	//fill holders
	vars := requestFiller.Fill(holders)
	fp.mergeVars(vars, fakeFiller.Fill(holders))

	m.Response.HttpHeaders, m.Response.Body = fp.walkAndFill(m.Response, vars)
	m.Callback.HttpHeaders, m.Callback.Body = fp.walkAndFill(m.Callback, vars)
}

func (fp ResponseMessageEvaluator) walkAndGet(res mock.ReplacementRequiredPayload) []string {
	vars := []string{}
	for _, header := range res.GetHeaders().Headers {
		for _, value := range header {
			fp.extractVars(value, &vars)
		}

	}

	for _, value := range res.GetHeaders().Cookies {
		fp.extractVars(value, &vars)
	}

	fp.extractVars(res.GetBody(), &vars)

	return vars
}

func (fp ResponseMessageEvaluator) walkAndFill(res mock.ReplacementRequiredPayload, vars map[string][]string) (mock.HttpHeaders, string) {
	headers := res.GetHeaders().Headers
	cookies := res.GetHeaders().Cookies
	body := res.GetBody()

	for header, values := range headers {
		for i, value := range values {
			headers[header][i] = fp.replaceVars(value, vars)
		}

	}

	for cookie, value := range cookies {
		cookies[cookie] = fp.replaceVars(value, vars)
	}

	newBody := fp.replaceVars(body, vars)

	return mock.HttpHeaders{Headers: headers, Cookies: cookies}, newBody
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
