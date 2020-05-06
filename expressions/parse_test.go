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
