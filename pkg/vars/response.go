package vars

import (
	"github.com/jmartin82/mmock/v3/pkg/mock"
	"strings"
)

type ResponseFiller struct {
	Response *mock.Response
}

func (rf ResponseFiller) Fill(holders []string) map[string][]string {
	vars := make(map[string][]string)
	log.Debugf("response: %v", rf.Response)
	bp := HttpEntityParams{Entity: &rf.Response.HTTPEntity}

	for _, tag := range holders {
		found := false
		s := ""

		if strings.HasPrefix(tag, "response.body.") {

			s, found = bp.getBodyParam(tag[14:])

		} else if strings.HasPrefix(tag, "response.cookie.") {

			s, found = bp.getCookieParam(tag[16:])

		} else if strings.HasPrefix(tag, "response.header.") {

			s, found = bp.getHeaderParam(tag[16:])
		}

		if found {
			vars[tag] = append(vars[tag], s)
		}

	}
	return vars
}
