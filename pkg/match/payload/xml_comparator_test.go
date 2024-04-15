package payload

import "testing"

func TestXMLComparator_Compare(t *testing.T) {
	type pathMap map[string]bool

	type args struct {
		s1            string
		s2            string
		optionalPaths pathMap
		currentPath   string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test same value order and format",
			args: args{
				s1:            "<note><to>Tove</to><from>Jani</from></note>",
				s2:            "<note><to>Tove</to><from>Jani</from></note>",
				optionalPaths: pathMap{},
				currentPath:   "",
			},
			want: true,
		},
		{
			name: "Test optionalPaths",
			args: args{
				s1: "<note><to><first>Tove</first><last>Boo</last></to><from>Jani</from></note>",
				s2: "<note><to><first>Tove</first></to><from>Jani</from></note>",
				optionalPaths: pathMap{
					".note.to.last": true,
				},
				currentPath: "",
			},
			want: true,
		},
		{
			name: "Test different order",
			args: args{
				s1:            "<note><to>Tove</to><from>Jani</from></note>",
				s2:            "<note><from>Jani</from><to>Tove</to></note>",
				optionalPaths: pathMap{},
				currentPath:   "",
			},
			want: true,
		},
		{
			name: "Test different format",
			args: args{
				s1:            "<note><to>Tove</to><from>Jani</from></note>",
				s2:            "<note> <to>Tove</to>\n<from>Jani</from></note>",
				optionalPaths: pathMap{},
				currentPath:   "",
			},
			want: true,
		},
		{
			name: "Test different values",
			args: args{
				s1:            "<?xml version=\"1.0\" encoding=\"utf-8\"?><note><to>Tove</to><from>Janid</from></note>",
				s2:            "<note><to>Tove</to><from>Jani</from></note>",
				optionalPaths: pathMap{},
				currentPath:   "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jc := &XMLComparator{}
			got := jc.Compare(tt.args.s1, tt.args.s2, tt.args.optionalPaths, tt.args.currentPath)
			if got != tt.want {
				t.Errorf("XMLComparator.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
