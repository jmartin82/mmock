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
			err:  true,
		}, {
			name: "float >1",
			in:   "1.2",
			err:  true,
		}, {
			name: "float <1",
			in:   "0.6",
		}, {
			name: "int",
			in:   "5",
			err:  true,
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
func TestDefinition_Validate(t *testing.T) {
	tests := []struct {
		name    string
		def     Definition
		wantErr bool
	}{
		{
			name: "valid definition",
			def: Definition{
				URI: "test-uri",
				Request: Request{
					Method: "GET",
					Path:   "/test-path",
				},
				Response: Response{
					StatusCode: 200,
				},
			},
			wantErr: false,
		},
		{
			name: "missing URI",
			def: Definition{
				Request: Request{
					Method: "GET",
					Path:   "/test-path",
				},
				Response: Response{
					StatusCode: 200,
				},
			},
			wantErr: true,
		},
		{
			name: "missing request method",
			def: Definition{
				URI: "test-uri",
				Request: Request{
					Path: "/test-path",
				},
				Response: Response{
					StatusCode: 200,
				},
			},
			wantErr: true,
		},
		{
			name: "missing request path",
			def: Definition{
				URI: "test-uri",
				Request: Request{
					Method: "GET",
				},
				Response: Response{
					StatusCode: 200,
				},
			},
			wantErr: true,
		},
		{
			name: "status code too low",
			def: Definition{
				URI: "test-uri",
				Request: Request{
					Method: "GET",
					Path:   "/test-path",
				},
				Response: Response{
					StatusCode: 99,
				},
			},
			wantErr: true,
		},
		{
			name: "status code too high",
			def: Definition{
				URI: "test-uri",
				Request: Request{
					Method: "GET",
					Path:   "/test-path",
				},
				Response: Response{
					StatusCode: 600,
				},
			},
			wantErr: true,
		},
		{
			name: "callback method without URL",
			def: Definition{
				URI: "test-uri",
				Request: Request{
					Method: "GET",
					Path:   "/test-path",
				},
				Response: Response{
					StatusCode: 200,
				},
				Callback: Callback{
					Method: "POST",
				},
			},
			wantErr: true,
		},
		{
			name: "callback URL without method",
			def: Definition{
				URI: "test-uri",
				Request: Request{
					Method: "GET",
					Path:   "/test-path",
				},
				Response: Response{
					StatusCode: 200,
				},
				Callback: Callback{
					Url: "http://example.com",
				},
			},
			wantErr: true,
		},
		{
			name: "valid with callback",
			def: Definition{
				URI: "test-uri",
				Request: Request{
					Method: "GET",
					Path:   "/test-path",
				},
				Response: Response{
					StatusCode: 200,
				},
				Callback: Callback{
					Method: "POST",
					Url:    "http://example.com",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.def.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Definition.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
