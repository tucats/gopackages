// Package expressions is a simple expression evaluator. It supports
// a rudementary symbol table with scoping, and knows about four data
// types (string, integer, double, and boolean). It does type casting as
// need automatically.
package expressions

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		want    interface{}
		wantErr bool
	}{
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
			name: "Simple addition",
			expr: "5 + i",
			want: 47,
		},
		{
			name: "Simple multiplication",
			expr: "5 * 20",
			want: 100,
		},
		{
			name: "Simple subtraction",
			expr: "pi - 1",
			want: 2.14,
		},
		{
			name: "Simple division",
			expr: "i / 7",
			want: 6,
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
			name: "Unary negation of single term",
			expr: "-i",
			want: -42,
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
			name: "len function",
			expr: "len(name) + 4",
			want: 7,
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
			want: nil,
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
			name: "Cast value to string",
			expr: "string(003.14)",
			want: "3.14",
		},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			e := New(tt.expr)
			symbols := map[string]interface{}{
				"i":    42,
				"pi":   3.14,
				"name": "Tom",
				"b":    true,
			}

			v1, err := e.Eval(symbols)
			if err != nil && tt.want != nil {
				t.Errorf("Expression test, unexpected error %v", err)
			} else {
				if v1 != tt.want {
					t.Errorf("Expression test, got %v, want %v", v1, tt.want)
				}
			}
		})
	}
}
