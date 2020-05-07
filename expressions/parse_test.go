package expressions

import (
	"reflect"
	"testing"
)

func TestExpression_Parse(t *testing.T) {
	type fields struct {
		Source string
		Type   ValueType
		Value  interface{}
		Tokens []string
		TokenP int
	}
	tests := []struct {
		name       string
		fields     fields
		wantErr    bool
		wantTokens []string
	}{
		{
			name: "Simple test",
			fields: fields{
				Source: "a = b",
			},
			wantErr:    false,
			wantTokens: []string{"a", "=", "b"},
		},
		{
			name: "Quote test",
			fields: fields{
				Source: "a = \"Tom\"",
			},
			wantErr:    false,
			wantTokens: []string{"a", "=", "\"Tom\""},
		},
		{
			name: "Compound operators",
			fields: fields{
				Source: "a >= \"Tom\"",
			},
			wantErr:    false,
			wantTokens: []string{"a", ">=", "\"Tom\""},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Expression{
				Source: tt.fields.Source,
				Type:   tt.fields.Type,
				Value:  tt.fields.Value,
				Tokens: tt.fields.Tokens,
				TokenP: tt.fields.TokenP,
			}
			if err := e.Parse(); (err != nil) != tt.wantErr {
				t.Errorf("Expression.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(e.Tokens, tt.wantTokens) {
				t.Errorf("Expression.Parse() got %v, want %v", e.Tokens, tt.wantTokens)
			}
		})
	}
}

func TestTokenize(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Simple alphanumeric name",
			args: args{
				src: "wage55",
			},
			want: []string{"wage55"},
		},
		{
			name: "Integer expression with spaces",
			args: args{
				src: "11 + 15",
			},
			want: []string{"11", "+", "15"},
		},
		{
			name: "Integer expression without spaces",
			args: args{
				src: "11+15",
			},
			want: []string{"11", "+", "15"},
		},
		{
			name: "String expression with spaces",
			args: args{
				src: "name + \"User\"",
			},
			want: []string{"name", "+", "\"User\""},
		},

		{
			name: "Float expression",
			args: args{
				src: "3.14 + 2",
			},
			want: []string{"3.14", "+", "2"},
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Tokenize(tt.args.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}
