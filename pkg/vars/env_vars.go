package vars

import (
	"os"
	"regexp"
	"strings"
)

var paramRe = regexp.MustCompile(`(?m)\((.*)\)`)

type EnvVars struct {
	valueMap map[string]string
	ready    bool
}

func (ev EnvVars) Fill(holders []string) map[string][]string {
	vars := make(map[string][]string)
	for _, tag := range holders {
		if strings.HasPrefix(tag, "env(") {
			var name = ev.getInputParam(tag)
			log.Debugf("looking for envar %s", name)
			vars[tag] = append(vars[tag], os.Getenv(name))
		}
	}
	return vars
}

func (ev EnvVars) getInputParam(param string) string {
	match := paramRe.FindStringSubmatch(param)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
