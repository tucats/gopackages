package util

import "testing"

func TestFormat(t *testing.T) {
	tests := []struct {
		name string
		arg  any
		want string
	}{
		{
			name: "Integer",
			arg:  33,
			want: "33",
		},
		{
			name: "Float",
			arg:  10.5,
			want: "10.5",
		},
		{
			name: "Array of int",
			arg:  []any{3, 5, 55},
			want: "[3, 5, 55]",
		},
		{
			name: "Array with array",
			arg:  []any{3, []any{"tom", true}, 55},
			want: "[3, [\"tom\", true], 55]",
		},
		{
			name: "simple structure",
			arg: map[string]any{
				"name": "Tom",
				"age":  59,
			},
			want: "{ age: 59, name: \"Tom\" }",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Format(tt.arg); got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}
