package payload

import "testing"

func TestJSONComparator_Compare(t *testing.T) {
	type args struct {
		s1 string
		s2 string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Test same value order and format", args{"{\"name\":\"bob\",\"age\":30}", "{\"name\":\"bob\",\"age\":30}"}, true},
		{"Test different order", args{"{\"name\":\"bob\",\"age\":30}", "{\"age\":30,\"name\":\"bob\"}"}, true},
		{"Test different format", args{"{\"name\":\"bob\",\"age\":30}", "{\"name\" : \"bob\"\n,\"age\" : 30}"}, true},
		{"Test different values", args{"{\"name\":\"bobs\",\"age\":30}", "{\"name\":\"bob\",\"age\":30}"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jc := &JSONComparator{}
			if got := jc.Compare(tt.args.s1, tt.args.s2); got != tt.want {
				t.Errorf("JSONComparator.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
