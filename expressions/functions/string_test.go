package functions

import (
	"reflect"
	"testing"
)

func TestFunctionLeft(t *testing.T) {
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
			name: "simple test",
			args: args{[]any{"Abraham", 4}},
			want: "Abra",
		},
		{
			name: "negative length test",
			args: args{[]any{"Abraham", -5}},
			want: "",
		},
		{
			name: "length too long test",
			args: args{[]any{"Abraham", 50}},
			want: "Abraham",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Left(nil, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("FunctionLeft() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FunctionLeft() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFunctionRight(t *testing.T) {
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
			name: "simple test",
			args: args{[]any{"Abraham", 3}},
			want: "ham",
		},
		{
			name: "length too small test",
			args: args{[]any{"Abraham", -5}},
			want: "",
		},
		{
			name: "length too long test",
			args: args{[]any{"Abraham", 103}},
			want: "Abraham",
		},
		{
			name: "empty string test",
			args: args{[]any{"", 3}},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Right(nil, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("FunctionRight() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FunctionRight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFunctionLower(t *testing.T) {
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
			name: "lower case",
			args: args{[]any{"short"}},
			want: "short",
		},
		{
			name: "upper case",
			args: args{[]any{"TALL"}},
			want: "tall",
		},
		{
			name: "mixed case",
			args: args{[]any{"camelCase"}},
			want: "camelcase",
		},
		{
			name: "empty string",
			args: args{[]any{""}},
			want: "",
		},
		{
			name: "non-string",
			args: args{[]any{3.14}},
			want: "3.14",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Lower(nil, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("FunctionLower() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FunctionLower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFunctionUpper(t *testing.T) {
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
			name: "lower case",
			args: args{[]any{"short"}},
			want: "SHORT",
		},
		{
			name: "upper case",
			args: args{[]any{"TALL"}},
			want: "TALL",
		},
		{
			name: "mixed case",
			args: args{[]any{"camelCase"}},
			want: "CAMELCASE",
		},
		{
			name: "empty string",
			args: args{[]any{""}},
			want: "",
		},
		{
			name: "non-string",
			args: args{[]any{3.14}},
			want: "3.14",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Upper(nil, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("FunctionUpper() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FunctionUpper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFunctionSubstring(t *testing.T) {
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
			name: "left case",
			args: args{[]any{"simple", 1, 3}},
			want: "sim",
		},
		{
			name: "right case",
			args: args{[]any{"simple", 3, 4}},
			want: "mple",
		},
		{
			name: "middle case",
			args: args{[]any{"simple", 3, 1}},
			want: "m",
		},
		{
			name: "invalid start case",
			args: args{[]any{"simple", -5, 3}},
			want: "sim",
		},
		{
			name: "invalid len case",
			args: args{[]any{"simple", 1, 355}},
			want: "simple",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Substring(nil, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("FunctionSubstring() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FunctionSubstring() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFunctionIndex(t *testing.T) {
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
			name: "index found",
			args: args{[]any{"string of text", "of"}},
			want: 8,
		},
		{
			name: "index not found",
			args: args{[]any{"string of text", "burp"}},
			want: 0,
		},
		{
			name: "empty source string",
			args: args{[]any{"", "burp"}},
			want: 0,
		},
		{
			name: "empty test string",
			args: args{[]any{"string of text", ""}},
			want: 1,
		},
		{
			name: "non-string test",
			args: args{[]any{"A1B2C3D4", 3}},
			want: 6,
		},
		{
			name: "array index",
			args: args{[]any{[]any{"tom", 3.14, true}, 3.14}},
			want: 2,
		},
		{
			name: "array not found",
			args: args{[]any{[]any{"tom", 3.14, true}, false}},
			want: 0,
		},
		{
			name: "member found",
			args: args{[]any{map[string]any{"name": "tom", "age": 55}, "age"}},
			want: true,
		},
		{
			name: "member found",
			args: args{[]any{map[string]any{"name": "tom", "age": 55}, "gender"}},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Index(nil, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("FunctionIndex() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FunctionIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}
