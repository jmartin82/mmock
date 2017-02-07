package match

import "github.com/jmartin82/mmock/definition"

type MatchVerifier struct {
	store   RequestStore
	matcher Matcher
}

func (mc MatchVerifier) Verify(r definition.Request) []definition.Request {
	requests := mc.store.GetRequests()
	matches := []definition.Request{}
	for _, req := range requests {
		if m, _ := mc.matcher.Match(&req, &definition.Mock{Request: r}, false); m {
			matches = append(matches, req)
		}
	}
	return matches

}

func NewMatchVerifier(matcher Matcher, requestStore RequestStore) *MatchVerifier {
	return &MatchVerifier{store: requestStore, matcher: matcher}
}
