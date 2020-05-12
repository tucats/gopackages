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

	for n := 0; n < b.emitPos; n++ {
		i := b.opcodes[n]
		opname, found := opcodeNames[i.Opcode]
		if !found {
			opname = fmt.Sprintf("Unknown %d", i.Opcode)
		}
		f := util.Format(i.Operand)
		if i.Operand == nil {
			f = ""
		}
		fmt.Printf("%4d: %s %s\n", n, opname, f)
	}

	fmt.Printf("\n%d instructions\n", b.emitPos)
}
