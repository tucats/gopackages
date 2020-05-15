package bytecode

import (
	"fmt"
	"strings"

	"github.com/tucats/gopackages/app-cli/ui"
	"github.com/tucats/gopackages/util"
)

var opcodeNames = map[int]string{
	Stop:               "Stop",
	Call:               "Call",
	Push:               "Push",
	Add:                "Add",
	Sub:                "Sub",
	Div:                "Div",
	Mul:                "Mul",
	And:                "And",
	Or:                 "Or",
	Negate:             "Negate",
	Branch:             "Branch",
	BranchTrue:         "BranchTrue",
	BranchFalse:        "BranchFalse",
	Equal:              "Equal",
	NotEqual:           "NotEqual",
	GreaterThan:        "GreaterThan",
	LessThan:           "LessThan",
	GreaterThanOrEqual: "GreaterThanOrEqual",
	LessThanOrEqual:    "LessThanOrEqual",
	Load:               "Load",
	Store:              "Store",
	Array:              "Array",
	LoadIndex:          "LoadIndex",
	StoreIndex:         "StoreIndex",
	Struct:             "Struct",
	Member:             "Member",
	Print:              "Print",
	Newline:            "Newline",
	Drop:               "Drop",
	AtLine:             "AtLine",
	MakeArray:          "MakeArray",
}

// Disasm prints out a representation of the bytecode for debugging purposes
func (b *ByteCode) Disasm() {

	// What is the maximum opcode name length?
	width := 0
	for _, k := range opcodeNames {
		if len(k) > width {
			width = len(k)
		}
	}

	ui.Debug("*** Disassembly %s", b.Name)
	for n := 0; n < b.emitPos; n++ {
		i := b.opcodes[n]
		opname, found := opcodeNames[i.Opcode]

		if !found {
			opname = fmt.Sprintf("Unknown %d", i.Opcode)
		}
		opname = (opname + strings.Repeat(" ", width))[:width]

		f := util.Format(i.Operand)
		if i.Operand == nil {
			f = ""
		}
		if i.Opcode >= BranchInstructions {
			f = "@" + f
		}
		ui.Debug("%4d: %s %s", n, opname, f)
	}

	ui.Debug("*** Disassembled %d instructions", b.emitPos)
}

// Format formats an array of bytecodes
func Format(opcodes []I) string {
	var b strings.Builder
	b.WriteRune('[')
	for n, i := range opcodes {

		if n > 0 {
			b.WriteRune(',')
		}
		opname, found := opcodeNames[i.Opcode]

		if !found {
			opname = fmt.Sprintf("Unknown %d", i.Opcode)
		}

		f := util.Format(i.Operand)
		if i.Operand == nil {
			f = ""
		}
		if i.Opcode >= BranchInstructions {
			f = "@" + f
		}
		b.WriteString(opname)
		if len(f) > 0 {
			b.WriteRune(' ')
			b.WriteString(f)
		}
	}
	b.WriteRune(']')
	return b.String()
}
