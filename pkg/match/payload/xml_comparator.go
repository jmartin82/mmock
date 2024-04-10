package payload

import (
	"encoding/json"
	"strings"

	xj "github.com/basgys/goxml2json"
)

type XMLComparator struct {
}

func (jc *XMLComparator) Compare(
  s1, s2 string,
  optionalPaths map[string]bool,
  currentPath string) bool {
	var o1 interface{}
	var o2 interface{}

	var err error

	b1, err := xj.Convert(strings.NewReader(s1))
	if err != nil {
		return false
	}

	b2, err := xj.Convert(strings.NewReader(s2))
	if err != nil {
		return false
	}

	err = json.Unmarshal(b1.Bytes(), &o1)
	if err != nil {
		return false
	}
	err = json.Unmarshal(b2.Bytes(), &o2)
	if err != nil {
		return false
	}
	
	json := &JSONComparator{}
	
	return json.doCompareJSONRegexUnmarshaled(
	  o1.(map[string]interface{}),
	  o2.(map[string]interface{}),
	  optionalPaths,
	  currentPath)
}
