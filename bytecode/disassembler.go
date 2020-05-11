package bytecode

import (
	"fmt"

	"github.com/tucats/gopackages/util"
)

var opcodeNames = map[int]string{
	Stop:               "Stop",
	Call:               "Call",
	Push:               "Push",
	Array:              "Array",
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
	Index:              "Index",
}

// Disasm prints out a representation of the bytecode for debugging purposes
func (b *ByteCode) Disasm() {

	for n, i := range b.opcodes {
		opname, found := opcodeNames[i.Opcode]
		if !found {
			opname = fmt.Sprintf("Unknown %d", i.Opcode)
		}
		fmt.Printf("%4d: %s %s\n", n, opname, util.Format(i.Operand))
	}

	fmt.Printf("\n%d instructions\n", len(b.opcodes))
}
