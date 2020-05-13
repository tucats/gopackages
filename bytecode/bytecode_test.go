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

func TestByteCode_Append(t *testing.T) {
	type fields struct {
		Name    string
		opcodes []I
		emitPos int
	}
	type args struct {
		a *ByteCode
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []I
		wantPos int
	}{
		{
			name: "simple append",
			fields: fields{
				opcodes: []I{
					I{Push, 0},
					I{Push, 0},
				},
				emitPos: 2,
			},
			args: args{
				a: &ByteCode{
					opcodes: []I{
						I{Add, nil},
					},
					emitPos: 1,
				},
			},
			want: []I{
				I{Push, 0},
				I{Push, 0},
				I{Add, nil},
			},
			wantPos: 3,
		},
		{
			name: "branch append",
			fields: fields{
				opcodes: []I{
					I{Push, 0},
					I{Push, 0},
				},
				emitPos: 2,
			},
			args: args{
				a: &ByteCode{
					opcodes: []I{
						I{Branch, 2}, // Must be updated
						I{Add, nil},
					},
					emitPos: 1,
				},
			},
			want: []I{
				I{Push, 0},
				I{Push, 0},
				I{Branch, 4}, // Updated from new offset
				I{Add, nil},
			},
			wantPos: 4,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ByteCode{
				Name:    tt.fields.Name,
				opcodes: tt.fields.opcodes,
				emitPos: tt.fields.emitPos,
			}
			b.Append(tt.args.a)
			if tt.wantPos != b.emitPos {
				t.Errorf("Append() wrong emitPos, got %d, want %d", b.emitPos, tt.wantPos)
			}
			// Check the slice of intentionally emitted opcodes (array may be larger)
			if !reflect.DeepEqual(tt.want, b.opcodes[:tt.wantPos]) {
				t.Errorf("Append() wrong array, got %v, want %v", b.opcodes, tt.want)
			}
		})
	}
}
