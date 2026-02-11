package payload

import (
	"strings"
)

type Comparer interface {
	Compare(s1, s2 string) bool
}

type Comparator struct {
	comparers map[string]Comparer
}

func NewComparator() *Comparator {
	comparers := make(map[string]Comparer)
	return &Comparator{comparers: comparers}
}

func NewDefaultComparator() *Comparator {
	comparator := NewComparator()
	json := &JSONComparator{}
	xml := &XMLComparator{}
	comparator.AddComparer("application/json", json)
	comparator.AddComparer("application/ld+json", json)
	comparator.AddComparer("application/merge-patch+json", json)
	comparator.AddComparer("application/xml", xml)
	comparator.AddComparer("text/xml", xml)

	// Register content type sniffer for generic binary types
	cts := NewContentTypeSniffer()
	comparator.AddComparer("application/octet-stream", cts)
	comparator.AddComparer("application/binary", cts)

	return comparator
}

func (c Comparator) AddComparer(contentType string, comparer Comparer) {
	c.comparers[contentType] = comparer
}

func (c Comparator) Compare(contentType, s1, s2 string) (comparable bool, equals bool) {
	parts := strings.Split(contentType, ";")
	comparer, ok := c.comparers[parts[0]]
	if !ok {
		return false, false
	}

	return true, comparer.Compare(s1, s2)
}

// SniffContentType detects the content type by inspecting the first non-whitespace character of the body.
func SniffContentType(body string) string {
	trimmedBody := strings.TrimLeft(body, " \t\r\n")
	if len(trimmedBody) > 0 {
		switch trimmedBody[0] {
		case '{', '[':
			return "application/json"
		case '<':
			return "application/xml"
		}
	}
	return ""
}

// ContentTypeSniffer implements Comparer by sniffing the body to detect its real type
// and delegating to the appropriate comparer.
type ContentTypeSniffer struct {
	comparers map[string]Comparer
}

func NewContentTypeSniffer() *ContentTypeSniffer {
	return &ContentTypeSniffer{
		comparers: map[string]Comparer{
			"application/json": &JSONComparator{},
			"application/xml":  &XMLComparator{},
		},
	}
}

func (cts *ContentTypeSniffer) Compare(s1, s2 string) bool {
	trimmed := strings.TrimLeft(s2, " \t\r\n")
	ct := SniffContentType(trimmed)
	if c, ok := cts.comparers[ct]; ok {
		return c.Compare(s1, trimmed)
	}
	return s1 == s2
}

// DetectedContentType returns the content type from headers, falling back to body sniffing if missing or generic.
func (cts *ContentTypeSniffer) DetectedContentType(headers map[string][]string, body string) string {
	if headers != nil {
		if ct, found := headers["Content-Type"]; found && len(ct) > 0 {
			val := ct[0]
			// For generic binary types, sniff the actual content
			if strings.HasPrefix(val, "application/octet-stream") || strings.HasPrefix(val, "application/binary") {
				if sniffed := SniffContentType(body); sniffed != "" {
					return sniffed
				}
			}
			return val
		}
	}
	return SniffContentType(body)
}

// IsJSON returns true if the detected content type is JSON.
func (cts *ContentTypeSniffer) IsJSON(headers map[string][]string, body string) bool {
	ct := cts.DetectedContentType(headers, body)
	return strings.HasPrefix(ct, "application/") && strings.HasSuffix(ct, "json")
}

// IsXML returns true if the detected content type is XML.
func (cts *ContentTypeSniffer) IsXML(headers map[string][]string, body string) bool {
	ct := cts.DetectedContentType(headers, body)
	return strings.HasPrefix(ct, "application/xml") || strings.HasPrefix(ct, "text/xml")
}

// IsFormEncoded returns true if the detected content type is x-www-form-urlencoded.
func (cts *ContentTypeSniffer) IsFormEncoded(headers map[string][]string, body string) bool {
	ct := cts.DetectedContentType(headers, body)
	return strings.HasPrefix(ct, "application/x-www-form-urlencoded")
}
