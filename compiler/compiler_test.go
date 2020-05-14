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
			name: "Simple print",
			arg:  "print 1",
			want: []bytecode.I{
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
				{Opcode: bytecode.Push, Operand: false},
				{Opcode: bytecode.Push, Operand: "bool"},
				{Opcode: bytecode.Call, Operand: 1},
				{Opcode: bytecode.BranchFalse, Operand: 8},
				{Opcode: bytecode.Push, Operand: 1},
				{Opcode: bytecode.Print, Operand: nil},
				{Opcode: bytecode.Newline, Operand: nil},
				{Opcode: bytecode.Branch, Operand: 11},
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
			bc, err := Compile(tokens)
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
