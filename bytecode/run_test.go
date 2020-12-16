package bytecode

import (
	"reflect"
	"testing"

	"github.com/tucats/gopackages/functions"
	"github.com/tucats/gopackages/symbols"
)

func TestByteCode_Run(t *testing.T) {
	type fields struct {
		Name    string
		opcodes []I
		emitPos int
		result  interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{

		{
			name: "stop",
			fields: fields{
				opcodes: []I{
					{Opcode: Stop},
				},
			},
		},
		{
			name: "push int",
			fields: fields{
				opcodes: []I{
					{Opcode: Push, Operand: 42},
					{Opcode: Stop},
				},
				result: 42,
			},
		},
		{
			name: "drop 2 stack items",
			fields: fields{
				opcodes: []I{
					{Opcode: Push, Operand: 42},
					{Opcode: Push, Operand: 43},
					{Opcode: Push, Operand: 44},
					{Opcode: Drop, Operand: 2},
					{Opcode: Stop},
				},
				result: 42,
			},
		},
		{
			name: "push float",
			fields: fields{
				opcodes: []I{
					{Opcode: Push, Operand: 3.14},
					{Opcode: Stop},
				},
				result: 3.14,
			},
		},
		{
			name: "add int",
			fields: fields{
				opcodes: []I{
					{Opcode: Push, Operand: 5},
					{Opcode: Push, Operand: 7},
					{Opcode: Add},
					{Opcode: Stop},
				},
				result: 12,
			},
		},
		{
			name: "add float to int",
			fields: fields{
				opcodes: []I{
					{Opcode: Push, Operand: 3.14},
					{Opcode: Push, Operand: 7},
					{Opcode: Add},
					{Opcode: Stop},
				},
				result: 10.14,
			},
		},
		{
			name: "sub int from  int",
			fields: fields{
				opcodes: []I{
					{Opcode: Push, Operand: 5},
					{Opcode: Push, Operand: 8},
					{Opcode: Sub},
					{Opcode: Stop},
				},
				result: -3,
			},
		},
		{
			name: "div float by int",
			fields: fields{
				opcodes: []I{
					{Opcode: Push, Operand: 10.0},
					{Opcode: Push, Operand: 2},
					{Opcode: Div},
					{Opcode: Stop},
				},
				result: 5.0,
			},
		},
		{
			name: "mul float by float",
			fields: fields{
				opcodes: []I{
					{Opcode: Push, Operand: 3.0},
					{Opcode: Push, Operand: 4.0},
					{Opcode: Mul},
					{Opcode: Stop},
				},
				result: 12.0,
			},
		},
		{
			name: "equal int test",
			fields: fields{
				opcodes: []I{
					{Opcode: Push, Operand: 5},
					{Opcode: Push, Operand: 5},
					{Opcode: Equal},
					{Opcode: Stop},
				},
				result: true,
			},
		},
		{
			name: "equal mixed test",
			fields: fields{
				opcodes: []I{
					{Opcode: Push, Operand: 5},
					{Opcode: Push, Operand: 5.0},
					{Opcode: Equal},
					{Opcode: Stop},
				},
				result: true,
			},
		},
		{
			name: "not equal int test",
			fields: fields{
				opcodes: []I{
					{Opcode: Push, Operand: 5},
					{Opcode: Push, Operand: 5},
					{Opcode: NotEqual},
					{Opcode: Stop},
				},
				result: false,
			},
		},
		{
			name: "not equal string test",
			fields: fields{
				opcodes: []I{
					{Opcode: Push, Operand: "fruit"},
					{Opcode: Push, Operand: "fruit"},
					{Opcode: NotEqual},
					{Opcode: Stop},
				},
				result: false,
			},
		},
		{
			name: "not equal bool test",
			fields: fields{
				opcodes: []I{
					{Opcode: Push, Operand: false},
					{Opcode: Push, Operand: false},
					{Opcode: NotEqual},
					{Opcode: Stop},
				},
				result: false,
			},
		},
		{
			name: "not equal float test",
			fields: fields{
				opcodes: []I{
					{Opcode: Push, Operand: 5.00000},
					{Opcode: Push, Operand: 5.00001},
					{Opcode: NotEqual},
					{Opcode: Stop},
				},
				result: true,
			},
		},
		{
			name: "greater than test",
			fields: fields{
				opcodes: []I{
					{Opcode: Push, Operand: 11.0},
					{Opcode: Push, Operand: 5},
					{Opcode: GreaterThan},
					{Opcode: Stop},
				},
				result: true,
			},
		},
		{
			name: "greater than or equals test",
			fields: fields{
				opcodes: []I{
					{Opcode: Push, Operand: "tom"},
					{Opcode: Push, Operand: "tom"},
					{Opcode: GreaterThanOrEqual},
					{Opcode: Stop},
				},
				result: true,
			},
		},
		{
			name: "length of string constant",
			fields: fields{
				opcodes: []I{
					{Opcode: Load, Operand: "len"},
					{Opcode: Push, Operand: "fruitcake"},
					{Opcode: Call, Operand: 1},
					{Opcode: Stop},
				},
				result: 9,
			},
		},
		{
			name: "left(n, 5) of string constant",
			fields: fields{
				opcodes: []I{
					// Arguments are pushed in the order parsed
					{Opcode: Load, Operand: "strings"},
					{Opcode: Push, Operand: "Left"},
					{Opcode: Member},
					{Opcode: Push, Operand: "fruitcake"},
					{Opcode: Push, Operand: 5},
					{Opcode: Call, Operand: 2},
					{Opcode: Stop},
				},
				result: "fruit",
			},
		},
		{
			name: "simple branch",
			fields: fields{
				opcodes: []I{
					{Opcode: Push, Operand: "fruitcake"},
					{Opcode: Branch, Operand: 3},
					{Opcode: Push, Operand: "Left"},
					{Opcode: Stop},
				},
				result: "fruitcake",
			},
		},
		{
			name: "if-true branch",
			fields: fields{
				opcodes: []I{

					// Use of "short-form" instruction initializer requires passing nil
					// for those instructions without an operand
					{Push, "stuff"},
					{Push, "fruitcake"},
					{Push, "fruitcake"},
					{Equal, nil},
					{BranchTrue, 6},
					{Push, .33333},
					{Stop, nil},
				},
				result: "stuff",
			},
		},
		{
			name: "if-false branch",
			fields: fields{
				opcodes: []I{

					// Use of "short-form" instruction initializer requires passing nil
					// for those instructions without an operand
					{Push, "stuff"},
					{Push, "cake"},
					{Push, "fruitcake"},
					{Equal, nil},
					{BranchTrue, 6},
					{Push, 42},
					{Stop, nil},
				},
				result: 42,
			},
		}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ByteCode{
				Name:    tt.fields.Name,
				opcodes: tt.fields.opcodes,
				emitPos: tt.fields.emitPos,
			}
			b.emitPos = len(b.opcodes)
			s := symbols.NewSymbolTable(tt.name)
			c := NewContext(s, b)
			functions.AddBuiltins(c.symbols)

			if err := c.Run(); (err != nil) != tt.wantErr {
				t.Errorf("ByteCode.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
			if c.running {
				t.Error("ByteCode Run() failed to stop interpreter")
			}
			if tt.fields.result != nil {
				v, err := c.Pop()
				if err != nil && !tt.wantErr {
					t.Error("ByteCode Run() unexpected " + err.Error())
				}
				if !reflect.DeepEqual(tt.fields.result, v) {
					t.Errorf("ByteCode Run() got %v, want %v ", v, tt.fields.result)

				}
			}
		})
	}
}
