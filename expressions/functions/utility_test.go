package functions

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
			name: "numeric value length",
			args: args{[]any{3.14}},
			want: 4,
		},
		{
			name: "array length",
			args: args{[]any{[]any{true, 3.14, "Tom"}}},
			want: 3,
		},
		{
			name: "struct value length",
			args: args{[]any{map[string]any{"name": "Tom", "age": 33}}},
			want: 2,
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

func TestFunctionProfile(t *testing.T) {
	type args struct {
		args []any
	}

	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{

		// Tests create an arbitrary key using a static UUID
		{
			name: "crete a key",
			args: args{[]any{"b306e250-6e07-4a05-abf4-e6a64d64cb72", "cookies"}},
			want: nil,
		},
		{
			name: "read a key",
			args: args{[]any{"b306e250-6e07-4a05-abf4-e6a64d64cb72"}},
			want: "cookies",
		},
		{
			name: "delete a key",
			args: args{[]any{"b306e250-6e07-4a05-abf4-e6a64d64cb72", ""}},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got any

			var err error
			if len(tt.args.args) > 1 {
				got, err = ProfileSet(nil, tt.args.args)
			} else {
				got, err = ProfileGet(nil, tt.args.args)
			}
			
			if (err != nil) != tt.wantErr {
				t.Errorf("FunctionProfile() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FunctionProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFunctionSort(t *testing.T) {
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
			name:    "scalar args",
			args:    args{[]any{66, 55}},
			want:    []any{55, 66},
			wantErr: false,
		},
		{
			name:    "mixed scalar args",
			args:    args{[]any{"tom", 3}},
			want:    []any{"3", "tom"},
			wantErr: false,
		},
		{
			name: "integer sort",
			args: args{[]any{[]any{55, 2, 18}}},
			want: []any{2, 18, 55},
		},
		{
			name: "float sort",
			args: args{[]any{[]any{55.0, 2, "18.5"}}},
			want: []any{2.0, 18.5, 55.0},
		},
		{
			name: "string sort",
			args: args{[]any{[]any{"pony", "cake", "unicorn", 5}}},
			want: []any{"5", "cake", "pony", "unicorn"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sort(nil, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("FunctionSort() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FunctionSort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFunctionMembers(t *testing.T) {
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
			name: "simple struct",
			args: args{[]any{map[string]any{"name": "Tom", "age": 55}}},
			want: []any{"age", "name"},
		},
		{
			name: "empty struct",
			args: args{[]any{map[string]any{}}},
			want: []any{},
		},
		{
			name:    "wrong type struct",
			args:    args{[]any{55}},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Members(nil, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("FunctionMembers() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FunctionMembers() = %v, want %v", got, tt.want)
			}
		})
	}
}
