package functions

import (
	"reflect"
	"testing"
)

func TestFunctionInt(t *testing.T) {
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
			name: "int(int)",
			args: args{[]any{33}},
			want: 33,
		},
		{
			name: "int(float64)",
			args: args{[]any{15.2}},
			want: 15,
		},
		{
			name: "int(string)",
			args: args{[]any{"42"}},
			want: 42,
		},
		{
			name: "int(bool)",
			args: args{[]any{true}},
			want: 1,
		},
		{
			name:    "int(error)",
			args:    args{[]any{"nescafe"}},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Int(nil, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("FunctionInt() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FunctionInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFunctionFloat(t *testing.T) {
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
			name: "float(int)",
			args: args{[]any{33}},
			want: 33.0,
		},
		{
			name: "float(float64)",
			args: args{[]any{15.2}},
			want: 15.2,
		},
		{
			name: "float(string)",
			args: args{[]any{"3.14"}},
			want: 3.14,
		},
		{
			name: "float(bool)",
			args: args{[]any{true}},
			want: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Float(nil, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("FunctionFloat() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FunctionFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFunctionString(t *testing.T) {
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
			name: "string(int)",
			args: args{[]any{33}},
			want: "33",
		},
		{
			name: "string(float64)",
			args: args{[]any{15.2}},
			want: "15.2",
		},
		{
			name: "string(string)",
			args: args{[]any{"3.14"}},
			want: "3.14",
		},
		{
			name: "string(bool)",
			args: args{[]any{true}},
			want: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := String(nil, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("FunctionString() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FunctionString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFunctionBool(t *testing.T) {
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
			name: "bool(int)",
			args: args{[]any{0}},
			want: false,
		},
		{
			name: "bool(float64)",
			args: args{[]any{15.2}},
			want: true,
		},
		{
			name: "bool(string)",
			args: args{[]any{"true"}},
			want: true,
		},
		{
			name: "bool(bool)",
			args: args{[]any{false}},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Bool(nil, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("FunctionBool() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FunctionBool() = %v, want %v", got, tt.want)
			}
		})
	}
}
