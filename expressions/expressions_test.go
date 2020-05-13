// Package expressions is a simple expression evaluator. It supports
// a rudementary symbol table with scoping, and knows about four data
// types (string, integer, double, and boolean). It does type casting as
// need automatically.
package expressions

import (
	"reflect"
	"testing"

	"github.com/tucats/gopackages/bytecode"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		want    interface{}
		wantErr bool
	}{
		{
			name: "Cast bool to string",
			expr: "string(b) + string(!b)",
			want: "truefalse",
		},
		{
			name: "Simple unary not",
			expr: "(!b)",
			want: false,
		},
		{
			name: "Simple division",
			expr: "i / 7",
			want: 6,
		},
		{
			name: "String subtraction",
			expr: "\"test\" - \"st\"",
			want: "te",
		},
		{
			name: "Case insensitive symbol names",
			expr: "name + NaMe",
			want: "TomTom",
		},
		{
			name: "Alphanumeric  symbol names",
			expr: "roman12 + \".\"",
			want: "XII.",
		},
		{
			name: "Simple addition",
			expr: "5 + i",
			want: 47,
		},
		{
			name: "Simple integer multiplication",
			expr: "5 * 20",
			want: 100,
		},
		{
			name: "Simple floating multiplication",
			expr: "5. * 20.",
			want: 100.,
		},
		{
			name: "Invalid type for  multiplication",
			expr: "true * \"many\"",
			want: nil,
		},
		{
			name: "Simple subtraction",
			expr: "pi - 1",
			want: 2.14,
		},
		{
			name: "Float division",
			expr: "10.0 / 4.0",
			want: 2.5,
		},
		{
			name: "Invalid division",
			expr: "\"house\" / \"cat\" ",
			want: nil,
		},
		{
			name: "Integer divide by zero ",
			expr: "17 / 0",
			want: nil,
		},
		{
			name: "Float divide by zero ",
			expr: "3.33 / 0.0 ",
			want: nil,
		},
		{
			name: "Order precedence",
			expr: "5 + i / 7",
			want: 11,
		},
		{
			name: "Order precedence with parens",
			expr: "(5 + i) * 2",
			want: 94,
		},
		{
			name: "Multiple paren terms",
			expr: "(i=42) & (name=\"Tom\")",
			want: true,
		},
		{
			name: "Invalid multiple paren terms",
			expr: "(i=42) & (name=\"Tom\"",
			want: nil,
		},
		{
			name: "Unary negation of single term",
			expr: "-i",
			want: -42,
		},
		{
			name: "Invalid unary negation",
			expr: "-name",
			want: nil,
		},
		{
			name: "Unary negation of diadic operator",
			expr: "43 + -i",
			want: 1,
		},
		{
			name: "Unary negation of subexpression",
			expr: "-(5+pi)",
			want: -8.14,
		},
		{
			name: "Type promotion bool to int",
			expr: "b + 3",
			want: 4,
		},
		{
			name: "Type promotion int to string",
			expr: "5 + name",
			want: "5Tom",
		},
		{
			name: "Type promotion int to float",
			expr: "pi + 5",
			want: 8.14,
		},
		{
			name: "Type coercion bool to int",
			expr: "int(true) + int(false)",
			want: 1,
		},
		{
			name: "Type coercion int to bool",
			expr: "i & true",
			want: true,
		},
		{
			name: "Type coercion string to bool",
			expr: "\"true\" | false",
			want: true,
		},
		{
			name: "Invalid type coercion string to bool",
			expr: "\"bob\" | false",
			want: false,
		},
		{
			name: "Cast value to int",
			expr: "int(3.14)",
			want: 3,
		},
		{
			name: "Cast value to float64",
			expr: "float(55)",
			want: 55.,
		},
		{
			name: "Cast bool to float64",
			expr: "float(b)",
			want: 1.,
		},
		{
			name: "Cast value to bool",
			expr: "bool(5)",
			want: true,
		},
		{
			name: "Cast float to string",
			expr: "string(003.14)",
			want: "3.14",
		},
		{
			name: "Cast int to string",
			expr: "string(i)",
			want: "42",
		},

		{
			name: "Homogeneous array constant",
			expr: "[1,2]",
			want: []interface{}{1, 2},
		},
		{
			name: "Heterogeneous array constant",
			expr: "[true,name, 33.5]",
			want: []interface{}{true, "Tom", 33.5},
		},
		{
			name: "Index into array",
			expr: "array[2]",
			want: "tom",
		},
		{
			name: "Add to array",
			expr: "[1,2] + 3",
			want: []interface{}{1, 2, 3},
		},
		{
			name: "Concatenate array",
			expr: "[1,2] + [3,4]",
			want: []interface{}{1, 2, 3, 4},
		},
		{
			name: "Subtract from array",
			expr: "[1,2,3] - 2",
			want: []interface{}{1, 3},
		}, {
			name: "Invalid argument list to function",
			expr: "len(1 3)",
			want: nil,
		},
		{
			name: "Incomplete argument list to function",
			expr: "len(13",
			want: nil,
		},
		{
			name: "Struct constant member",
			expr: "{name:\"Tom\", age:50}.age",
			want: 50,
		},
		{
			name: "Struct constant nested array member",
			expr: "{name:\"Tom\",sizes:[10, 12], age:50}.sizes",
			want: []interface{}{10, 12},
		},
		{
			name: "Struct constant nested array member indexed",
			expr: "{name:\"Tom\",sizes:[10, 12], age:50}.sizes[2]",
			want: 12,
		}, {
			name: "len of string function",
			expr: "len(name) + 4",
			want: 7,
		},
		{
			name: "len of array function",
			expr: "len(array) + 4",
			want: 8,
		},
		{
			name: "left function",
			expr: "left(name, 2)",
			want: "To",
		},
		{
			name: "right function",
			expr: "right(name, 2)",
			want: "om",
		},
		{
			name: "index of string function",
			expr: "index(name, \"o\")",
			want: 2,
		},
		{
			name: "index not found function",
			expr: "index(name, \"g\")",
			want: 0,
		},
		{
			name: "index of array function",
			expr: "index(array, false)",
			want: 4,
		},
		{
			name: "index of array not found",
			expr: "index(array, 55.5)",
			want: 0,
		},
		{
			name: "substring function",
			expr: "substring(\"ABCDEF\", 2, 3)",
			want: "BCD",
		},
		{
			name: "empty substring function",
			expr: "substring(\"ABCDEF\", 5, 0)",
			want: "",
		},
		{
			name: "Invalid argument count to function",
			expr: "substring(\"ABCDEF\", 5)",
			want: nil,
		},
		{
			name: "upper function",
			expr: "upper(name)",
			want: "TOM",
		},
		{
			name: "lower function",
			expr: "lower(name)",
			want: "tom",
		},
		{
			name: "min homogeneous args function",
			expr: "min(15,33,11,6)",
			want: 6,
		},
		{
			name: "min float args function",
			expr: "min(3.0, 1.0, 2.0)",
			want: 1.0,
		},
		{
			name: "min string args function",
			expr: "min(\"house\", \"cake\", \"pig\" )",
			want: "cake",
		},
		{
			name: "min hetergenous args function",
			expr: "min(15,33.5,\"11\",6)",
			want: 6,
		},
		{
			name: "max hetergenous args function",
			expr: "max(15.1,33.5,\"11\",6)",
			want: 33.5,
		},
		{
			name: "sum hetergenous args function",
			expr: "sum(10.1, 5, \"2\")",
			want: 17.1,
		},
		{
			name: "sum homogeneous args function",
			expr: "sum(\"abc\", \"137\", \"def\")",
			want: "abc137def",
		},
		{
			name: "Simple true conditional",
			expr: "true ? 1.1 : 1.0",
			want: 1.1,
		},
		{
			name: "Simple false conditional",
			expr: "5 = 6 ? 1.1 : 1.0",
			want: 1.0,
		},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			e := New(tt.expr)

			// Create a common symbol table.
			symbols := bytecode.NewSymbolTable(tt.name)
			symbols.Set("i", 42)
			symbols.Set("pi", 3.14)
			symbols.Set("name", "Tom")
			symbols.Set("b", true)
			symbols.Set("roman12", "XII")
			symbols.Set("array", []interface{}{1, "tom", 33., false})

			v1, err := e.Eval(symbols)
			if err != nil && tt.want != nil {
				t.Errorf("Expression test, unexpected error %v", err)
			} else {
				if !reflect.DeepEqual(v1, tt.want) {
					t.Errorf("Expression test, got %v, want %v", v1, tt.want)
				}
			}
		})
	}
}
