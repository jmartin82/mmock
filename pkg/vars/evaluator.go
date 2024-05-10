package vars

import (
	"regexp"
	"strings"

	"github.com/jmartin82/mmock/v3/pkg/mock"
	"github.com/jmartin82/mmock/v3/pkg/match"
)

var varsRegex = regexp.MustCompile(`\{\{\s*(.+?)\s*\}\}`)

type Evaluator interface {
	Eval(req *mock.Request, m *mock.Definition, scenearioStore match.ScenearioStorer)
}

type ResponseMessageEvaluator struct {
	FillerFactory FillerFactory
}

func NewResponseMessageEvaluator(fp FillerFactory) *ResponseMessageEvaluator {
	return &ResponseMessageEvaluator{FillerFactory: fp}
}

func (fp ResponseMessageEvaluator) Eval(req *mock.Request, m *mock.Definition, store match.ScenearioStorer) {
	requestFiller := fp.FillerFactory.CreateRequestFiller(req, m)
	fakeFiller := fp.FillerFactory.CreateFakeFiller()
	streamFiller := fp.FillerFactory.CreateStreamFiller()
	responseFiller := fp.FillerFactory.CreateResponseFiller(&m.Response)

	//first replace the external streams
	holders := fp.walkAndGet(m.Response.HTTPEntity)
	holders = append(holders, fp.walkAndGet(m.Callback.HTTPEntity)...)

	//fill holders with the correct values
	vars := streamFiller.Fill(holders)
	fp.walkAndFill(&m.Response.HTTPEntity, vars)
	fp.walkAndFill(&m.Callback.HTTPEntity, vars)

	//repeat the same opration in order to replace any holder
	//coming from the external streams

	//get the holders in the response and the callback structs
	holders = fp.walkAndGet(m.Response.HTTPEntity)
	holders = append(holders, fp.walkAndGet(m.Callback.HTTPEntity)...)

	//fill holders with the correct values
	vars = requestFiller.Fill(holders)
	fp.mergeVars(vars, fakeFiller.Fill(holders))

	// if we have a scenario, fill any scenario.* holders
	if m.Control.Scenario.Name != "" {
	  scenarioFiller := fp.FillerFactory.CreateScenarioFiller(req, m, store, m.Control.Scenario.Name)
	  fp.mergeVars(vars, scenarioFiller.Fill(holders))
	}

	//replace the holders in the response and the callback
	fp.walkAndFill(&m.Response.HTTPEntity, vars)

	// fill any response.* holders
	fp.mergeVars(vars, responseFiller.Fill(holders))

	//replace the holders in the response
	fp.walkAndFill(&m.Callback.HTTPEntity, vars)
}

func (fp ResponseMessageEvaluator) walkAndGet(res mock.HTTPEntity) []string {

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

func (fp ResponseMessageEvaluator) walkAndFill(res *mock.HTTPEntity, vars map[string][]string) {
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
