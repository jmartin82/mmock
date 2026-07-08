package payload

import "testing"

func TestComparator_Compare(t *testing.T) {
	type fields struct {
		comparers map[string]Comparer
	}
	type args struct {
		contentType string
		s1          string
		s2          string
	}

	tests := []struct {
		name           string
		fields         fields
		args           args
		wantComparable bool
		wantEquals     bool
	}{
		{"Compare json ok", fields{map[string]Comparer{"application/json": &JSONComparator{}}}, args{"application/json", "{\"name\":\"bob\",\"age\":30}", "{\"name\":\"bob\",\"age\":30}"}, true, true},
		{"Compare json ko", fields{map[string]Comparer{"application/json": &JSONComparator{}}}, args{"application/json", "{\"name\":\"bob\",\"age\":30}", "{\"name\":\"bob\",\"age\":40}"}, true, false},
		{"Not comparable", fields{map[string]Comparer{"application/xml": &XMLComparator{}}}, args{"application/json", "{\"name\":\"bob\",\"age\":30}", "{\"name\":\"bob\",\"age\":40}"}, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Comparator{
				comparers: tt.fields.comparers,
			}
			gotComparable, gotEquals := c.Compare(tt.args.contentType, tt.args.s1, tt.args.s2)
			if gotComparable != tt.wantComparable {
				t.Errorf("Comparator.Compare() gotComparable = %v, want %v", gotComparable, tt.wantComparable)
			}
			if gotEquals != tt.wantEquals {
				t.Errorf("Comparator.Compare() gotEquals = %v, want %v", gotEquals, tt.wantEquals)
			}
		})
	}
}

func TestContentTypeSniffer_DetectedContentType(t *testing.T) {
	jsonBody := `{"name":"bob"}`
	xmlBody := `<root><name>bob</name></root>`

	tests := []struct {
		name    string
		headers map[string][]string
		body    string
		want    string
	}{
		{"missing header, json body", nil, jsonBody, "application/json"},
		{"missing header, xml body", map[string][]string{}, xmlBody, "application/xml"},
		{"json header trusted", map[string][]string{"Content-Type": {"application/json"}}, jsonBody, "application/json"},
		{"json header with charset stripped", map[string][]string{"Content-Type": {"application/json; charset=utf-8"}}, jsonBody, "application/json"},
		{"aws json header trusted", map[string][]string{"Content-Type": {"application/x-amz-json-1.1"}}, jsonBody, "application/x-amz-json-1.1"},
		{"octet-stream sniffed to json", map[string][]string{"Content-Type": {"application/octet-stream"}}, jsonBody, "application/json"},
		{"binary sniffed to xml", map[string][]string{"Content-Type": {"application/binary"}}, xmlBody, "application/xml"},
		{"wrong type application/text sniffed to json", map[string][]string{"Content-Type": {"application/text"}}, jsonBody, "application/json"},
		{"wrong type text/plain sniffed to xml", map[string][]string{"Content-Type": {"text/plain"}}, xmlBody, "application/xml"},
		{"wrong type with unstructured body keeps header", map[string][]string{"Content-Type": {"application/text"}}, "hello", "application/text"},
	}

	cts := NewContentTypeSniffer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cts.DetectedContentType(tt.headers, tt.body); got != tt.want {
				t.Errorf("DetectedContentType() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestContentTypeSniffer_IsJSON(t *testing.T) {
	jsonBody := `{"name":"bob"}`
	xmlBody := `<root><name>bob</name></root>`

	tests := []struct {
		name    string
		headers map[string][]string
		body    string
		want    bool
	}{
		{"application/json", map[string][]string{"Content-Type": {"application/json"}}, jsonBody, true},
		{"application/ld+json", map[string][]string{"Content-Type": {"application/ld+json"}}, jsonBody, true},
		{"application/json with charset", map[string][]string{"Content-Type": {"application/json; charset=utf-8"}}, jsonBody, true},
		{"application/x-amz-json-1.1", map[string][]string{"Content-Type": {"application/x-amz-json-1.1"}}, jsonBody, true},
		{"missing header json body", nil, jsonBody, true},
		{"octet-stream json body", map[string][]string{"Content-Type": {"application/octet-stream"}}, jsonBody, true},
		{"wrong type application/text json body", map[string][]string{"Content-Type": {"application/text"}}, jsonBody, true},
		{"application/xml", map[string][]string{"Content-Type": {"application/xml"}}, xmlBody, false},
		{"wrong type text/plain xml body", map[string][]string{"Content-Type": {"text/plain"}}, xmlBody, false},
		{"application/text unstructured body", map[string][]string{"Content-Type": {"application/text"}}, "hello", false},
	}

	cts := NewContentTypeSniffer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cts.IsJSON(tt.headers, tt.body); got != tt.want {
				t.Errorf("IsJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComparatorDefaultFactory(t *testing.T) {
	c := NewDefaultComparator()

	comparers := []string{
		"application/json",
		"application/ld+json",
		"application/merge-patch+json",
		"application/xml",
		"text/xml",
		"application/octet-stream",
		"application/binary",
	}

	for _, comparer := range comparers {
		if _, ok := c.comparers[comparer]; !ok {
			t.Errorf("%s content type doesn't have comparator", comparer)
		}
	}
}
