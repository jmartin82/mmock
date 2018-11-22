package vars

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`(?m)\((.*)\)`)

type Stream struct {
}

func (st Stream) Fill(holders []string) map[string][]string {

	vars := make(map[string][]string)
	for _, tag := range holders {
		if strings.HasPrefix(tag, "file.contents(") {
			vars[tag] = append(vars[tag], st.getOutput(st.getFileContents(tag)))
		} else if strings.HasPrefix(tag, "http.contents(") {
			vars[tag] = append(vars[tag], st.getOutput(st.getHttpContents(tag)))
		}
	}
	return vars
}

func (st Stream) getOutput(o []byte, err error) string {

	if err != nil {
		log.Printf("Impossible read mock stream: %s", err)
		return fmt.Sprintf("ERROR: %s", err.Error())
	}
	return string(o)
}

func (st Stream) getFileContents(tag string) ([]byte, error) {
	path := st.getInputParam(tag)
	return ioutil.ReadFile(path)
}

func (st Stream) getHttpContents(tag string) ([]byte, error) {
	url := st.getInputParam(tag)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (st Stream) getInputParam(param string) string {
	match := re.FindStringSubmatch(param)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
