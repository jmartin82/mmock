package mock

import (
	"testing"
	"time"
)

func TestJSONParseDelay(t *testing.T) {
	tests := []struct {
		name string
		in   string
		exp  time.Duration
		err  bool
	}{
		{
			name: "null",
			in:   "null",
			err:  true,
		}, {
			name: "empty string",
			in:   `""`,
			err:  true,
		}, {
			name: "object",
			in:   "{}",
			err:  true,
		}, {
			name: "array",
			in:   "[]",
			err:  true,
		}, {
			name: "float 1",
			in:   "1.00",
			exp:  time.Second,
		}, {
			name: "float >1",
			in:   "1.2",
			exp:  time.Second,
		}, {
			name: "float <1",
			in:   "0.6",
		}, {
			name: "int",
			in:   "5",
			exp:  time.Second * 5,
		}, {
			name: "string",
			in:   `"5"`,
			err:  true,
		}, {
			name: "string",
			in:   `"5s"`,
			exp:  time.Second * 5,
		}, {
			name: "valid",
			in:   `"1m40s"`,
			exp:  time.Minute + time.Second*40,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Delay{}
			err := d.UnmarshalJSON([]byte(tt.in))
			if tt.err && err == nil {
				t.Errorf("expected error: got: %v", d.Duration)
			}
			if want, got := tt.exp, d.Duration; want != got {
				t.Errorf("want: %v, got: %v", want, got)
			}
		})
	}
}
