package jsonrpc

import "testing"

func Test_isJSONArray(t *testing.T) {
	type args struct {
		d []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty",
			args: args{
				d: []byte(``),
			},
			want: false,
		},
		{
			name: "object",
			args: args{
				d: []byte(`{}`),
			},
			want: false,
		},
		{
			name: "array",
			args: args{
				d: []byte(`[]`),
			},
			want: true,
		},
		{
			name: "array with spaces",
			args: args{
				d: []byte(`   		
[ 		
]   				
`),
			},
			want: true,
		},
		{
			name: "json",
			args: args{
				d: []byte(`{"a":null,"b":"f"}`),
			},
			want: false,
		},
		{
			name: "json 2",
			args: args{
				d: []byte(`[{"a":null,"b":"f"},{"a":null,"b":"f"}]`),
			},
			want: true,
		},
		{
			name: "complex json",
			args: args{
				d: []byte(`[{"a":"b"},{"a":{"b":1,"c":["a", 1, {"b":2}]}]`),
			},
			want: true,
		},
		{
			name: "complex json 2",
			args: args{
				d: []byte(`{"a":{"b":1,"c":["a", 1, {"b":2}]}`),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isJSONArray(tt.args.d); got != tt.want {
				t.Errorf("isJSONArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
