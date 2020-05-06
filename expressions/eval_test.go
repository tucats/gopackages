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
