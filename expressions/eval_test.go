package expressions

import (
	"reflect"
	"testing"
)

func TestExpression_Eval(t *testing.T) {
	type fields struct {
		Source   string
		Type     ValueType
		Value    interface{}
		Tokens   []string
		TokenPos []int
		TokenP   int
	}
	type args struct {
		symbols map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Order precedence",
			fields: fields{
				Tokens: []string{"5", "+", "2", "*", "3"},
			},
			args: args{
				symbols: map[string]interface{}{"a": 42, "pi": 3.14},
			},
			want:    11,
			wantErr: false,
		},
		{
			name: "Order precedence - parens",
			fields: fields{
				Tokens: []string{"(", "5", "+", "2", ")", "*", "3"},
			},
			args: args{
				symbols: map[string]interface{}{"a": 42, "pi": 3.14},
			},
			want:    21,
			wantErr: false,
		},
		{
			name: "Mixed type string+int addition",
			fields: fields{
				Tokens: []string{"\"user\"", "+", "a"},
			},
			args: args{
				symbols: map[string]interface{}{"a": 42, "pi": 3.14},
			},
			want:    "user42",
			wantErr: false,
		},
		{
			name: "Simple integer",
			fields: fields{
				Tokens: []string{"1"},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "Simple float",
			fields: fields{
				Tokens: []string{"1.2"},
			},
			want:    1.2,
			wantErr: false,
		},
		{
			name: "Simple string",
			fields: fields{
				Tokens: []string{"\"Test\""},
			},
			want:    "Test",
			wantErr: false,
		},
		{
			name: "Simple bool",
			fields: fields{
				Tokens: []string{"true"},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Simple variable",
			fields: fields{
				Tokens: []string{"pi"},
			},
			args: args{
				symbols: map[string]interface{}{"a": 42, "pi": 3.14},
			},
			want:    3.14,
			wantErr: false,
		},
		{
			name: "Alphanumeric variable",
			fields: fields{
				Tokens: []string{"rate22"},
			},
			args: args{
				symbols: map[string]interface{}{"a": 42, "rate22": 15.25},
			},
			want:    15.25,
			wantErr: false,
		},
		{
			name: "Invalid variable",
			fields: fields{
				Tokens: []string{"cost"},
			},
			args: args{
				symbols: map[string]interface{}{"a": 42, "pi": 3.14},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Simple integer addition",
			fields: fields{
				Tokens: []string{"3", "+", "4"},
			},
			args: args{
				symbols: map[string]interface{}{"a": 42, "pi": 3.14},
			},
			want:    7,
			wantErr: false,
		},
		{
			name: "Simple integer subtraction",
			fields: fields{
				Tokens: []string{"11", "-", "2"},
			},
			args: args{
				symbols: map[string]interface{}{"a": 42, "pi": 3.14},
			},
			want:    9,
			wantErr: false,
		},
		{
			name: "Simple integer const+variable",
			fields: fields{
				Tokens: []string{"11", "+", "a"},
			},
			args: args{
				symbols: map[string]interface{}{"a": 42, "pi": 3.14},
			},
			want:    53,
			wantErr: false,
		},
		{
			name: "Mixed type float+int addition",
			fields: fields{
				Tokens: []string{"3.14", "+", "a"},
			},
			args: args{
				symbols: map[string]interface{}{"a": 42, "pi": 3.14},
			},
			want:    45.14,
			wantErr: false,
		},
		{
			name: "Mixed type bool+int addition",
			fields: fields{
				Tokens: []string{"true", "+", "a"},
			},
			args: args{
				symbols: map[string]interface{}{"a": 42, "pi": 3.14},
			},
			want:    43,
			wantErr: false,
		},
		{
			name: "Integer equals",
			fields: fields{
				Tokens: []string{"42", "=", "a"},
			},
			args: args{
				symbols: map[string]interface{}{"a": 42, "pi": 3.14},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "String equals true",
			fields: fields{
				Tokens: []string{"name", "=", "\"tom\""},
			},
			args: args{
				symbols: map[string]interface{}{"name": "tom"},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "String equals false",
			fields: fields{
				Tokens: []string{"name", "=", "\"tim\""},
			},
			args: args{
				symbols: map[string]interface{}{"name": "tom"},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Integer GE true",
			fields: fields{
				Tokens: []string{"55", ">=", "33"},
			},
			args: args{
				symbols: map[string]interface{}{"name": "tom"},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Boolean GE invalid",
			fields: fields{
				Tokens: []string{"true", ">=", "a"},
			},
			args: args{
				symbols: map[string]interface{}{"a": false},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Boolean type coercion",
			fields: fields{
				Tokens: []string{"1.0", "|", "true"},
			},
			args: args{
				symbols: map[string]interface{}{"a": false},
			},
			want:    true,
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Expression{
				Source:   tt.fields.Source,
				Type:     tt.fields.Type,
				Value:    tt.fields.Value,
				Tokens:   tt.fields.Tokens,
				TokenPos: tt.fields.TokenPos,
				TokenP:   tt.fields.TokenP,
			}
			got, err := e.Eval(tt.args.symbols)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expression.Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expression.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}
