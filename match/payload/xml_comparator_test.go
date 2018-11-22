package payload

import "testing"

func TestXMLComparator_Compare(t *testing.T) {
	type args struct {
		s1 string
		s2 string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Test same value order and format", args{"<note><to>Tove</to><from>Jani</from></note>", "<note><to>Tove</to><from>Jani</from></note>"}, true},
		{"Test different order", args{"<note><to>Tove</to><from>Jani</from></note>", "<note><from>Jani</from><to>Tove</to></note>"}, true},
		{"Test different format", args{"<note><to>Tove</to><from>Jani</from></note>", "<note> <to>Tove</to>\n<from>Jani</from></note>"}, true},
		{"Test different values", args{"<?xml version=\"1.0\" encoding=\"utf-8\"?><note><to>Tove</to><from>Janid</from></note>", "<note><to>Tove</to><from>Jani</from></note>"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jc := &XMLComparator{}
			if got := jc.Compare(tt.args.s1, tt.args.s2); got != tt.want {
				t.Errorf("XMLComparator.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
