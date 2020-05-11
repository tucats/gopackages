package bytecode

import (
	"reflect"
	"testing"
)

func TestByteCode_Run(t *testing.T) {
	type fields struct {
		Name    string
		opcodes []I
		emitPos int
		pc      int
		stack   []interface{}
		sp      int
		running bool
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
					I{opcode: Stop},
				},
			},
		},
		{
			name: "push int",
			fields: fields{
				opcodes: []I{
					I{opcode: Push, operand: 42},
					I{opcode: Stop},
				},
				result: 42,
			},
		},
		{
			name: "push float",
			fields: fields{
				opcodes: []I{
					I{opcode: Push, operand: 3.14},
					I{opcode: Stop},
				},
				result: 3.14,
			},
		},
		{
			name: "add int",
			fields: fields{
				opcodes: []I{
					I{opcode: Push, operand: 5},
					I{opcode: Push, operand: 7},
					I{opcode: Add},
					I{opcode: Stop},
				},
				result: 12,
			},
		},
		{
			name: "add float to int",
			fields: fields{
				opcodes: []I{
					I{opcode: Push, operand: 3.14},
					I{opcode: Push, operand: 7},
					I{opcode: Add},
					I{opcode: Stop},
				},
				result: 10.14,
			},
		},
		{
			name: "sub int from  int",
			fields: fields{
				opcodes: []I{
					I{opcode: Push, operand: 8},
					I{opcode: Push, operand: 5},
					I{opcode: Sub},
					I{opcode: Stop},
				},
				result: -3,
			},
		},
		{
			name: "div float by int",
			fields: fields{
				opcodes: []I{
					I{opcode: Push, operand: 2},
					I{opcode: Push, operand: 10.0},
					I{opcode: Div},
					I{opcode: Stop},
				},
				result: 5.0,
			},
		},
		{
			name: "mul float by float",
			fields: fields{
				opcodes: []I{
					I{opcode: Push, operand: 3.0},
					I{opcode: Push, operand: 4.0},
					I{opcode: Mul},
					I{opcode: Stop},
				},
				result: 12.0,
			},
		},
		{
			name: "equal int test",
			fields: fields{
				opcodes: []I{
					I{opcode: Push, operand: 5},
					I{opcode: Push, operand: 5},
					I{opcode: Equal},
					I{opcode: Stop},
				},
				result: true,
			},
		},
		{
			name: "equal mixed test",
			fields: fields{
				opcodes: []I{
					I{opcode: Push, operand: 5},
					I{opcode: Push, operand: 5.0},
					I{opcode: Equal},
					I{opcode: Stop},
				},
				result: true,
			},
		},
		{
			name: "not equal int test",
			fields: fields{
				opcodes: []I{
					I{opcode: Push, operand: 5},
					I{opcode: Push, operand: 5},
					I{opcode: NotEqual},
					I{opcode: Stop},
				},
				result: false,
			},
		},
		{
			name: "not equal string test",
			fields: fields{
				opcodes: []I{
					I{opcode: Push, operand: "fruit"},
					I{opcode: Push, operand: "fruit"},
					I{opcode: NotEqual},
					I{opcode: Stop},
				},
				result: false,
			},
		},
		{
			name: "not equal bool test",
			fields: fields{
				opcodes: []I{
					I{opcode: Push, operand: false},
					I{opcode: Push, operand: false},
					I{opcode: NotEqual},
					I{opcode: Stop},
				},
				result: false,
			},
		},
		{
			name: "not equal float test",
			fields: fields{
				opcodes: []I{
					I{opcode: Push, operand: 5.00000},
					I{opcode: Push, operand: 5.00001},
					I{opcode: NotEqual},
					I{opcode: Stop},
				},
				result: true,
			},
		},
		{
			name: "greater than test",
			fields: fields{
				opcodes: []I{
					I{opcode: Push, operand: 11.0},
					I{opcode: Push, operand: 5},
					I{opcode: GreaterThan},
					I{opcode: Stop},
				},
				result: true,
			},
		},
		{
			name: "greater than or equals test",
			fields: fields{
				opcodes: []I{
					I{opcode: Push, operand: "tom"},
					I{opcode: Push, operand: "tom"},
					I{opcode: GreaterThanOrEqual},
					I{opcode: Stop},
				},
				result: true,
			},
		},
		{
			name: "length of string constant",
			fields: fields{
				opcodes: []I{
					I{opcode: Push, operand: "fruitcake"},
					I{opcode: Push, operand: "len"},
					I{opcode: Call, operand: 1},
					I{opcode: Stop},
				},
				result: 9,
			},
		},
		{
			name: "left(n, 5) of string constant",
			fields: fields{
				opcodes: []I{
					// Arguments are pushed in the order parsed
					I{opcode: Push, operand: "fruitcake"},
					I{opcode: Push, operand: 5},
					I{opcode: Push, operand: "left"},
					I{opcode: Call, operand: 2},
					I{opcode: Stop},
				},
				result: "fruit",
			},
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ByteCode{
				Name:    tt.fields.Name,
				opcodes: tt.fields.opcodes,
				emitPos: tt.fields.emitPos,
				pc:      tt.fields.pc,
				stack:   tt.fields.stack,
				sp:      tt.fields.sp,
				running: tt.fields.running,
			}
			if err := b.Run(); (err != nil) != tt.wantErr {
				t.Errorf("ByteCode.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
			if b.running {
				t.Error("ByteCode Run() failed to stop interpreter")
			}
			if tt.fields.result != nil {
				v, err := b.Pop()
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
