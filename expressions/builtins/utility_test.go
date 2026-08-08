package builtins

import (
	"reflect"
	"testing"
)

func TestFunctionLen(t *testing.T) {
	type args struct {
		args []any
	}

	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "string length",
			args: args{[]any{"hamster"}},
			want: 7,
		},
		{
			name: "empty string length",
			args: args{[]any{""}},
			want: 0,
		},
		{
			name:    "numeric value length",
			args:    args{[]any{3.14}},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Length(nil, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("FunctionLen() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FunctionLen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLength(t *testing.T) {
	tests := []struct {
		name    string
		args    []any
		want    any
		wantErr bool
	}{
		{
			name: "simple string",
			args: []any{"foo"},
			want: 3,
		},
		{
			name: "unicode string",
			args: []any{"\u2318foo\u2318"},
			want: 9,
		},
		{
			name: "int converted to string",
			args: []any{"123456"},
			want: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Length(nil, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Length() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Length() = %v, want %v", got, tt.want)
			}
		})
	}
}
