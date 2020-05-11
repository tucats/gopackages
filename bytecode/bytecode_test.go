package bytecode

import (
	"reflect"
	"testing"
)

func TestByteCode_New(t *testing.T) {
	t.Run("test", func(t *testing.T) {

		b := New("testing")

		want := ByteCode{
			Name:    "testing",
			opcodes: make([]I, InitialOpcodeSize),
			emitPos: 0,
			pc:      0,
			sp:      0,
			stack:   make([]interface{}, InitialStackSize),
			running: false,
		}
		if !reflect.DeepEqual(*b, want) {
			t.Error("new() did not return expected object")
		}
	})
}

func TestByteCode_Emit(t *testing.T) {
	type fields struct {
		Name    string
		opcodes []I
		emitPos int
		pc      int
		stack   []interface{}
		sp      int
		running bool
		emit    []I
	}
	type args struct {
		emit []I
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []I
	}{
		{
			name: "first emit",
			fields: fields{
				opcodes: []I{},
			},
			args: args{
				emit: []I{
					I{Push, 33},
				},
			},
			want: []I{
				I{Push, 33},
			},
		},
		{
			name: "multiple emit",
			fields: fields{
				opcodes: []I{},
			},
			args: args{
				emit: []I{
					I{Push, 33},
					I{Push, "stuff"},
					I{Opcode: Add},
					I{Opcode: Stop},
				},
			},
			want: []I{
				I{Push, 33},
				I{Push, "stuff"},
				I{Opcode: Add},
				I{Opcode: Stop},
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
			for _, i := range tt.args.emit {
				b.Emit(i.Opcode, i.Operand)
			}

			for n, i := range b.opcodes {
				if n < b.emitPos && i != tt.want[n] {
					t.Error("opcode mismatch")
				}
			}
		})
	}
}
