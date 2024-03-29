package builtins

import (
	"reflect"
	"testing"

	"github.com/tucats/gopackages/defs"
	"github.com/tucats/gopackages/expressions/symbols"
)

func TestFunctionIndex(t *testing.T) {
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
			name: "index found",
			args: args{[]interface{}{"string of text", "of"}},
			want: 8,
		},
		{
			name: "index not found",
			args: args{[]interface{}{"string of text", "burp"}},
			want: 0,
		},
		{
			name: "empty source string",
			args: args{[]interface{}{"", "burp"}},
			want: 0,
		},
		{
			name: "empty test string",
			args: args{[]interface{}{"string of text", ""}},
			want: 1,
		},
		{
			name: "non-string test",
			args: args{[]interface{}{"A1B2C3D4", 3}},
			want: 6,
		},
		{
			name: "array index",
			args: args{[]interface{}{[]interface{}{"tom", 3.14, true}, 3.14}},
			want: 1,
		},
		{
			name: "array not found",
			args: args{[]interface{}{[]interface{}{"tom", 3.14, true}, false}},
			want: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We will need a symbol table so the Index function can find out
			// if it is allowed or not.
			s := symbols.NewSymbolTable("testing")
			s.Root().SetAlways(defs.ExtensionsVariable, true)

			got, err := Index(s, tt.args.args)
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
