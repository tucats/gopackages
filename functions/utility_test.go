package functions

import (
	"reflect"
	"testing"
)

func TestFunctionLen(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "string length",
			args: args{[]interface{}{"hamster"}},
			want: 7,
		},
		{
			name: "empty string length",
			args: args{[]interface{}{""}},
			want: 0,
		},
		{
			name: "numeric value length",
			args: args{[]interface{}{3.14}},
			want: 4,
		},
		{
			name: "array length",
			args: args{[]interface{}{[]interface{}{true, 3.14, "Tom"}}},
			want: 3,
		},
		{
			name: "struct value length",
			args: args{[]interface{}{map[string]interface{}{"name": "Tom", "age": 33}}},
			want: 2,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FunctionLen(nil, tt.args.args)
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
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{

		// Tests create an arbitrary key using a static UUID
		{
			name: "crete a key",
			args: args{[]interface{}{"b306e250-6e07-4a05-abf4-e6a64d64cb72", "cookies"}},
			want: true,
		},
		{
			name: "read a key",
			args: args{[]interface{}{"b306e250-6e07-4a05-abf4-e6a64d64cb72"}},
			want: "cookies",
		},
		{
			name: "delete a key",
			args: args{[]interface{}{"b306e250-6e07-4a05-abf4-e6a64d64cb72", ""}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FunctionProfile(nil, tt.args.args)
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
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:    "bad arg type",
			args:    args{[]interface{}{55}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "integer sort",
			args: args{[]interface{}{[]interface{}{55, 2, 18}}},
			want: []interface{}{2, 18, 55},
		},
		{
			name: "float sort",
			args: args{[]interface{}{[]interface{}{55.0, 2, "18.5"}}},
			want: []interface{}{2.0, 18.5, 55.0},
		},
		{
			name: "string sort",
			args: args{[]interface{}{[]interface{}{"pony", "cake", "unicorn", 5}}},
			want: []interface{}{"5", "cake", "pony", "unicorn"},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FunctionSort(nil, tt.args.args)
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
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "simple struct",
			args: args{[]interface{}{map[string]interface{}{"name": "Tom", "age": 55}}},
			want: []interface{}{"age", "name"},
		},
		{
			name: "empty struct",
			args: args{[]interface{}{map[string]interface{}{}}},
			want: []interface{}{},
		},
		{
			name:    "wrong type struct",
			args:    args{[]interface{}{55}},
			want:    nil,
			wantErr: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FunctionMembers(nil, tt.args.args)
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
