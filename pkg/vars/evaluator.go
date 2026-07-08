package vars

import (
	"regexp"
	"strings"

	"github.com/jmartin82/mmock/v3/pkg/mock"
)

var varsRegex = regexp.MustCompile(`\{\{\s*(.+?)\s*\}\}`)

// streamStartRegex matches only the opening of a stream tag. The whole tag is then
// delimited by a parenthesis-depth scan (findStreamSpan) instead of a regex, so the
// argument may contain nested {{...}} holders and gjson parentheses/braces.
var streamStartRegex = regexp.MustCompile(`\{\{\s*(?:file|http)\.contents\(`)

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
	streamFiller := fp.FillerFactory.CreateStreamFiller(req, m)

	//first replace the external streams
	holders := fp.walkAndGetStreams(m.Response.HTTPEntity)
	vars := streamFiller.Fill(holders)
	fp.walkAndFillStreams(&m.Response.HTTPEntity, vars)

	//repeat the same opration in order to replace any holder coming from the external streams

	//get the holders in the response and the callback structs
	holders = fp.walkAndGet(m.Response.HTTPEntity)
	holders = append(holders, fp.walkAndGet(m.Callback.HTTPEntity)...)

	//fill holders with the correct values
	vars = requestFiller.Fill(holders)
	fp.mergeVars(vars, fakeFiller.Fill(holders))

	//replace the holders in the response and the callback
	fp.walkAndFill(&m.Response.HTTPEntity, vars)
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

// walkAndGetStreams collects the inner content of every whole stream tag
// ({{ file.contents(...) }} / {{ http.contents(...) }}) found in the entity, using
// findStreamSpan so arguments containing nested {{...}} holders are captured intact.
func (fp ResponseMessageEvaluator) walkAndGetStreams(res mock.HTTPEntity) []string {
	vars := []string{}
	for _, header := range res.Headers {
		for _, value := range header {
			vars = append(vars, extractStreamTags(value)...)
		}
	}
	for _, value := range res.Cookies {
		vars = append(vars, extractStreamTags(value)...)
	}
	vars = append(vars, extractStreamTags(res.Body)...)
	return vars
}

// walkAndFillStreams replaces every whole stream tag with its resolved value.
func (fp ResponseMessageEvaluator) walkAndFillStreams(res *mock.HTTPEntity, vars map[string][]string) {
	for header, values := range res.Headers {
		for i, value := range values {
			res.Headers[header][i] = replaceStreamTags(value, vars)
		}
	}
	for cookie, value := range res.Cookies {
		res.Cookies[cookie] = replaceStreamTags(value, vars)
	}
	res.Body = replaceStreamTags(res.Body, vars)
}

// findStreamSpan locates the next whole stream tag at or after `from`, returning the
// byte offsets of the whole {{ ... }} span and the trimmed inner content (e.g.
// "file.contents(...)"). Unlike varsRegex, it tolerates nested {{...}} holders and
// gjson parentheses/braces by scanning for the balanced closing ')' of contents(...)
// before requiring the closing '}}'. Malformed openings are skipped.
func findStreamSpan(s string, from int) (start, end int, inner string, ok bool) {
	for from < len(s) {
		loc := streamStartRegex.FindStringIndex(s[from:])
		if loc == nil {
			return 0, 0, "", false
		}
		start = from + loc[0]
		openParen := from + loc[1] // index just after the '(' of contents(

		depth := 1
		i := openParen
		for ; i < len(s); i++ {
			if s[i] == '(' {
				depth++
			} else if s[i] == ')' {
				depth--
				if depth == 0 {
					break
				}
			}
		}

		if depth == 0 {
			j := i + 1
			for j < len(s) && (s[j] == ' ' || s[j] == '\t' || s[j] == '\r' || s[j] == '\n') {
				j++
			}
			if j+1 < len(s) && s[j] == '}' && s[j+1] == '}' {
				return start, j + 2, strings.TrimSpace(s[start+2 : i+1]), true
			}
		}

		from = openParen // malformed opening: skip past it and keep searching
	}
	return 0, 0, "", false
}

// extractStreamTags returns the inner content of every whole stream tag in input.
func extractStreamTags(input string) []string {
	tags := []string{}
	for from := 0; ; {
		_, end, inner, ok := findStreamSpan(input, from)
		if !ok {
			break
		}
		tags = append(tags, inner)
		from = end
	}
	return tags
}

// replaceStreamTags substitutes each whole stream tag with the first unused value from
// vars keyed by the tag's inner content, leaving unresolved tags verbatim. It uses the
// same one-shot per-value semantics as replaceVars.
func replaceStreamTags(input string, vars map[string][]string) string {
	var b strings.Builder
	from := 0
	for {
		start, end, inner, ok := findStreamSpan(input, from)
		if !ok {
			break
		}
		b.WriteString(input[from:start])
		if v, found := vars[inner]; found && len(v) > 0 {
			b.WriteString(v[0])
			vars[inner] = v[1:]
		} else {
			b.WriteString(input[start:end])
		}
		from = end
	}
	b.WriteString(input[from:])
	return b.String()
}
