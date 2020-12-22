package compiler

import (
	"reflect"
	"testing"

	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/tokenizer"
)

func TestCompile(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    []bytecode.I
		wantErr bool
	}{
		{
			name: "for index loop with break",
			arg:  "for i := 0; i < 10; i = i + 1 { if i == 3 { break }; print i }",
			want: []bytecode.I{
				{Opcode: bytecode.AtLine, Operand: 1},
				{Opcode: bytecode.PushScope, Operand: nil},
				{Opcode: bytecode.Push, Operand: 0},
				{Opcode: bytecode.SymbolCreate, Operand: "i"},
				{Opcode: bytecode.Store, Operand: "i"},
				{Opcode: bytecode.Load, Operand: "i"},
				{Opcode: bytecode.Push, Operand: 10},
				{Opcode: bytecode.LessThan, Operand: nil},
				{Opcode: bytecode.BranchFalse, Operand: 33},
				{Opcode: bytecode.AtLine, Operand: 1},
				{Opcode: bytecode.PushScope, Operand: nil},
				{Opcode: bytecode.AtLine, Operand: 1},
				{Opcode: bytecode.Load, Operand: "bool"},
				{Opcode: bytecode.Load, Operand: "i"},
				{Opcode: bytecode.Push, Operand: 3},
				{Opcode: bytecode.Equal, Operand: nil},
				{Opcode: bytecode.Call, Operand: 1},
				{Opcode: bytecode.BranchFalse, Operand: 23},
				{Opcode: bytecode.AtLine, Operand: 1},
				{Opcode: bytecode.PushScope, Operand: nil},
				{Opcode: bytecode.AtLine, Operand: 1},
				{Opcode: bytecode.Branch, Operand: 33},
				{Opcode: bytecode.PopScope, Operand: nil},
				{Opcode: bytecode.AtLine, Operand: 1},
				{Opcode: bytecode.Load, Operand: "i"},
				{Opcode: bytecode.Print, Operand: nil},
				{Opcode: bytecode.Newline, Operand: nil},
				{Opcode: bytecode.PopScope, Operand: nil},
				{Opcode: bytecode.Load, Operand: "i"},
				{Opcode: bytecode.Push, Operand: 1},
				{Opcode: bytecode.Add, Operand: nil},
				{Opcode: bytecode.Store, Operand: "i"},
				{Opcode: bytecode.Branch, Operand: 5},
				{Opcode: bytecode.PopScope, Operand: nil},
			},
		},
		{
			name: "Simple block",
			arg:  "{ print ; print } ",
			want: []bytecode.I{
				{Opcode: bytecode.AtLine, Operand: 1},
				{Opcode: bytecode.PushScope, Operand: nil},
				{Opcode: bytecode.AtLine, Operand: 1},
				{Opcode: bytecode.Newline, Operand: nil},
				{Opcode: bytecode.AtLine, Operand: 1},
				{Opcode: bytecode.Newline, Operand: nil},
				{Opcode: bytecode.PopScope, Operand: nil},
			},
			wantErr: false,
		},
		{
			name: "store to _",
			arg:  "_ = 3",
			want: []bytecode.I{
				{Opcode: bytecode.AtLine, Operand: 1},
				{Opcode: bytecode.Push, Operand: 3},
				{Opcode: bytecode.Drop, Operand: 1},
			},
			wantErr: false,
		},
		{
			name: "Simple print",
			arg:  "print 1",
			want: []bytecode.I{
				{Opcode: bytecode.AtLine, Operand: 1},
				{Opcode: bytecode.Push, Operand: 1},
				{Opcode: bytecode.Print, Operand: nil},
				{Opcode: bytecode.Newline, Operand: nil},
			},
			wantErr: false,
		},
		{
			name: "Simple if else",
			arg:  "if false print 1 else print 2",
			want: []bytecode.I{
				{Opcode: bytecode.AtLine, Operand: 1},
				{Opcode: bytecode.Load, Operand: "bool"},
				{Opcode: bytecode.Push, Operand: false},
				{Opcode: bytecode.Call, Operand: 1},
				{Opcode: bytecode.BranchFalse, Operand: 10},
				{Opcode: bytecode.AtLine, Operand: 1},
				{Opcode: bytecode.Push, Operand: 1},
				{Opcode: bytecode.Print, Operand: nil},
				{Opcode: bytecode.Newline, Operand: nil},
				{Opcode: bytecode.Branch, Operand: 14},
				{Opcode: bytecode.AtLine, Operand: 1},
				{Opcode: bytecode.Push, Operand: 2},
				{Opcode: bytecode.Print, Operand: nil},
				{Opcode: bytecode.Newline, Operand: nil},
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := tokenizer.New(tt.arg)
			c := New()
			// Make sure PRINT verb works for these tests.
			c.printEnabled = true
			bc, err := c.Compile(tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opcodes := bc.Opcodes()
			if !reflect.DeepEqual(opcodes, tt.want) {
				t.Errorf("Compile() = %v, want %v", bytecode.Format(opcodes), bytecode.Format(tt.want))
			}
		})
	}
}
