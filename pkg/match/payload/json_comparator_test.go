package payload

import "testing"

func TestJSONComparator_Compare(t *testing.T) {

	type pathMap map[string]bool

	type args struct {
		s1 string
		s2 string
		optionalPaths pathMap
		currentPath string
	}

	tests := []struct {
		name string
		args args
		want bool
	      }{

		{
		  name: "Test same value order and format", 
		  args: args{
		    s1: "{\"name\":\"bob\",\"age\":30}",
		    s2: "{\"name\":\"bob\",\"age\":30}",
		    optionalPaths: pathMap{},
		    currentPath: "",
		  },
		  want: true},
		  {
		    name: "Test same value order and format in regex", 
		    args: args{
		      s1: "{\"name\":\"b.*\",\"age\":30}",
		      s2: "{\"name\":\"bob\",\"age\":30}",
		      optionalPaths: pathMap{},
		      currentPath: "",
		    },
		    want: true,
		  },
		  {
		    name: "Test different order", 
		    args: args{
		      s1: "{\"name\":\"bob\",\"age\":30}",
		      s2: "{\"age\":30,\"name\":\"bob\"}",
		      optionalPaths: pathMap{},
		      currentPath: "",
		    },
		    want: true,
		  },
		  {
		    name: "Test equal arrays", 
		    args: args{
		      s1: "[{\"name\":\"bob\",\"age\":30}]",
		      s2: "[{\"age\":30,\"name\":\"bob\"}]",
		      optionalPaths: pathMap{},
		      currentPath: "",
		    },
		    want: true,
		  },
		  {
		    name: "Test object and array difference", 
		    args: args{
		      s1: "{\"name\":\"bob\",\"age\":30}",
		      s2: "[{\"age\":30,\"name\":\"bob\"}]",
		      optionalPaths: pathMap{},
		      currentPath: "",
		    },
		    want: false,
		  },
		  {
		    name: "Test array with object with regex", 
		    args: args{
		      s1: "[{\"name\":\".*\",\"age\":\"\\\\d*\"}]",
		      s2: "[{\"age\":30,\"name\":\"bob\"}]",
		      optionalPaths: pathMap{},
		      currentPath: ".",
		    },
		    want: true,
		  },
		  {
		    name: "Test array of objects with regex", 
		    args: args{
		      s1: "[{\"name\":\".*\",\"age\":\"\\\\d*\"}]",
		      s2: "[{\"age\":3,\"name\":\"gary\"},{\"age\":30,\"name\":\"bob\"}]",
		      optionalPaths: pathMap{},
		      currentPath: "",
		    },
		    want: true,
		  },
		  {
		    name: "Test nested objects with regex", 
		    args: args{
		      s1: "{\"name\": {\"firstName\":\".*\"},\"age\":\"\\\\d*\"}",
		      s2: "{\"age\":30,\"name\": {\"firstName\": \"bob\"}}",
		      optionalPaths: pathMap{},
		      currentPath: "",
		    },
		    want: true,
		  },
		  {
		    name: "Test nested objects with regex with optionalPaths",
		    args: args{
		      s1: "{\"name\": {\"firstName\":\".*\", \"lastName\":\".*\"},\"age\":\"\\\\d*\"}",
		      s2: "{\"age\":30,\"name\": {\"firstName\": \"bob\"}}",
		      optionalPaths: pathMap{".name.lastName": true,},
		      currentPath: "",
		    },
		    want: true,
		  },
		  {
		    name: "Test array of nested objects with regex", 
		    args: args{
		      s1: "[{\"name\": {\"firstName\":\".*\"},\"age\":\"\\\\d*\"}]",
		      s2: "[{\"age\":3,\"name\": {\"firstName\": \"gary\"}},{\"age\":30,\"name\": {\"firstName\": \"bob\"}}]",
		      optionalPaths: pathMap{},
		      currentPath: "",
		    },
		    want: true,
		  },
		  {
		    name: "Test different arrays", 
		    args: args{
		      s1: "[{\"name\":\"john\",\"age\":30}]",
		      s2: "[{\"age\":30,\"name\":\"bob\"}]",
		      optionalPaths: pathMap{},
		      currentPath: "",
		    },
		    want: false,
		  },
		  {
		    name: "Test different format", 
		    args: args{
		      s1: "{\"name\":\"bob\",\"age\":30}",
		      s2: "{\"name\" : \"bob\"\n,\"age\" : 30}",
		      optionalPaths: pathMap{},
		      currentPath: "",
		    },
		    want: true,
		  },
		  {
		    name: "Test different values", 
		    args: args{
		      s1: "{\"name\":\"bobs\",\"age\":30}",
		      s2: "{\"name\":\"bob\",\"age\":30}",
		      optionalPaths: pathMap{},
		      currentPath: "",
		    },
		    want: false,
		  },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jc := &JSONComparator{}
			if got := jc.Compare(tt.args.s1, tt.args.s2, tt.args.optionalPaths, tt.args.currentPath); got != tt.want {
				t.Errorf("JSONComparator.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
