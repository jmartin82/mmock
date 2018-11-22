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

func TestComparatorDefaultFactory(t *testing.T) {

	c := NewDefaultComparator()

	if _, ok := c.comparers["application/json"]; !ok {
		t.Errorf("application/json content type doesn't have comparator")
	}
	if _, ok := c.comparers["application/xml"]; !ok {
		t.Errorf("application/json content type doesn't have comparator")
	}
	if _, ok := c.comparers["text/xml"]; !ok {
		t.Errorf("application/json content type doesn't have comparator")
	}
}
