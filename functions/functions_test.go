package functions

import (
	"reflect"
	"testing"
)

func TestFunctionInt(t *testing.T) {
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
			name: "int(int)",
			args: args{[]interface{}{33}},
			want: 33,
		},
		{
			name: "int(float64)",
			args: args{[]interface{}{15.2}},
			want: 15,
		},
		{
			name: "int(string)",
			args: args{[]interface{}{"42"}},
			want: 42,
		},
		{
			name: "int(bool)",
			args: args{[]interface{}{true}},
			want: 1,
		},
		{
			name:    "int(error)",
			args:    args{[]interface{}{"nescafe"}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FunctionInt(tt.args.args)
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
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "float(int)",
			args: args{[]interface{}{33}},
			want: 33.0,
		},
		{
			name: "float(float64)",
			args: args{[]interface{}{15.2}},
			want: 15.2,
		},
		{
			name: "float(string)",
			args: args{[]interface{}{"3.14"}},
			want: 3.14,
		},
		{
			name: "float(bool)",
			args: args{[]interface{}{true}},
			want: 1.0,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FunctionFloat(tt.args.args)
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

func TestFunctionLeft(t *testing.T) {
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
			name: "simple test",
			args: args{[]interface{}{"Abraham", 4}},
			want: "Abra",
		},
		{
			name: "negative length test",
			args: args{[]interface{}{"Abraham", -5}},
			want: "",
		},
		{
			name: "length too long test",
			args: args{[]interface{}{"Abraham", 50}},
			want: "Abraham",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FunctionLeft(tt.args.args)
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
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "simple test",
			args: args{[]interface{}{"Abraham", 3}},
			want: "ham",
		},
		{
			name: "length too small test",
			args: args{[]interface{}{"Abraham", -5}},
			want: "",
		},
		{
			name: "length too long test",
			args: args{[]interface{}{"Abraham", 103}},
			want: "Abraham",
		},
		{
			name: "empty string test",
			args: args{[]interface{}{"", 3}},
			want: "",
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FunctionRight(tt.args.args)
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
