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
		symbols map[string]interface{}
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
					I{Opcode: Stop},
				},
			},
		},
		{
			name: "push int",
			fields: fields{
				opcodes: []I{
					I{Opcode: Push, Operand: 42},
					I{Opcode: Stop},
				},
				result: 42,
			},
		},
		{
			name: "push float",
			fields: fields{
				opcodes: []I{
					I{Opcode: Push, Operand: 3.14},
					I{Opcode: Stop},
				},
				result: 3.14,
			},
		},
		{
			name: "add int",
			fields: fields{
				opcodes: []I{
					I{Opcode: Push, Operand: 5},
					I{Opcode: Push, Operand: 7},
					I{Opcode: Add},
					I{Opcode: Stop},
				},
				result: 12,
			},
		},
		{
			name: "add float to int",
			fields: fields{
				opcodes: []I{
					I{Opcode: Push, Operand: 3.14},
					I{Opcode: Push, Operand: 7},
					I{Opcode: Add},
					I{Opcode: Stop},
				},
				result: 10.14,
			},
		},
		{
			name: "sub int from  int",
			fields: fields{
				opcodes: []I{
					I{Opcode: Push, Operand: 5},
					I{Opcode: Push, Operand: 8},
					I{Opcode: Sub},
					I{Opcode: Stop},
				},
				result: -3,
			},
		},
		{
			name: "div float by int",
			fields: fields{
				opcodes: []I{
					I{Opcode: Push, Operand: 10.0},
					I{Opcode: Push, Operand: 2},
					I{Opcode: Div},
					I{Opcode: Stop},
				},
				result: 5.0,
			},
		},
		{
			name: "mul float by float",
			fields: fields{
				opcodes: []I{
					I{Opcode: Push, Operand: 3.0},
					I{Opcode: Push, Operand: 4.0},
					I{Opcode: Mul},
					I{Opcode: Stop},
				},
				result: 12.0,
			},
		},
		{
			name: "equal int test",
			fields: fields{
				opcodes: []I{
					I{Opcode: Push, Operand: 5},
					I{Opcode: Push, Operand: 5},
					I{Opcode: Equal},
					I{Opcode: Stop},
				},
				result: true,
			},
		},
		{
			name: "equal mixed test",
			fields: fields{
				opcodes: []I{
					I{Opcode: Push, Operand: 5},
					I{Opcode: Push, Operand: 5.0},
					I{Opcode: Equal},
					I{Opcode: Stop},
				},
				result: true,
			},
		},
		{
			name: "not equal int test",
			fields: fields{
				opcodes: []I{
					I{Opcode: Push, Operand: 5},
					I{Opcode: Push, Operand: 5},
					I{Opcode: NotEqual},
					I{Opcode: Stop},
				},
				result: false,
			},
		},
		{
			name: "not equal string test",
			fields: fields{
				opcodes: []I{
					I{Opcode: Push, Operand: "fruit"},
					I{Opcode: Push, Operand: "fruit"},
					I{Opcode: NotEqual},
					I{Opcode: Stop},
				},
				result: false,
			},
		},
		{
			name: "not equal bool test",
			fields: fields{
				opcodes: []I{
					I{Opcode: Push, Operand: false},
					I{Opcode: Push, Operand: false},
					I{Opcode: NotEqual},
					I{Opcode: Stop},
				},
				result: false,
			},
		},
		{
			name: "not equal float test",
			fields: fields{
				opcodes: []I{
					I{Opcode: Push, Operand: 5.00000},
					I{Opcode: Push, Operand: 5.00001},
					I{Opcode: NotEqual},
					I{Opcode: Stop},
				},
				result: true,
			},
		},
		{
			name: "greater than test",
			fields: fields{
				opcodes: []I{
					I{Opcode: Push, Operand: 11.0},
					I{Opcode: Push, Operand: 5},
					I{Opcode: GreaterThan},
					I{Opcode: Stop},
				},
				result: true,
			},
		},
		{
			name: "greater than or equals test",
			fields: fields{
				opcodes: []I{
					I{Opcode: Push, Operand: "tom"},
					I{Opcode: Push, Operand: "tom"},
					I{Opcode: GreaterThanOrEqual},
					I{Opcode: Stop},
				},
				result: true,
			},
		},
		{
			name: "length of string constant",
			fields: fields{
				opcodes: []I{
					I{Opcode: Push, Operand: "fruitcake"},
					I{Opcode: Push, Operand: "len"},
					I{Opcode: Call, Operand: 1},
					I{Opcode: Stop},
				},
				result: 9,
			},
		},
		{
			name: "left(n, 5) of string constant",
			fields: fields{
				opcodes: []I{
					// Arguments are pushed in the order parsed
					I{Opcode: Push, Operand: "fruitcake"},
					I{Opcode: Push, Operand: 5},
					I{Opcode: Push, Operand: "left"},
					I{Opcode: Call, Operand: 2},
					I{Opcode: Stop},
				},
				result: "fruit",
			},
		},
		{
			name: "simple branch",
			fields: fields{
				opcodes: []I{
					I{Opcode: Push, Operand: "fruitcake"},
					I{Opcode: Branch, Operand: 3},
					I{Opcode: Push, Operand: "left"},
					I{Opcode: Stop},
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
					I{Push, "stuff"},
					I{Push, "fruitcake"},
					I{Push, "fruitcake"},
					I{Equal, nil},
					I{BranchTrue, 6},
					I{Push, .33333},
					I{Stop, nil},
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
					I{Push, "stuff"},
					I{Push, "cake"},
					I{Push, "fruitcake"},
					I{Equal, nil},
					I{BranchTrue, 6},
					I{Push, 42},
					I{Stop, nil},
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
				pc:      tt.fields.pc,
				stack:   tt.fields.stack,
				sp:      tt.fields.sp,
				running: tt.fields.running,
			}
			b.emitPos = len(b.opcodes)
			if err := b.Run(map[string]interface{}{}); (err != nil) != tt.wantErr {
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
