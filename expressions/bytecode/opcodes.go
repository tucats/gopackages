package bytecode

/*
 * ADDING A NEW OPCODE
 *
 * 1. Add the Opcode name as a constant in the list below. If it is an opcode
 *    that has a bytecode address as its operand, put it in the section
 *    identified as "branch instructions".
 *
 * 2. Add the opcode name to the map below, which converts the const identifier
 *    to a human-readable name. By convention, the human-readable name is the same as
 *    the constant itself.
 *
 * 3. Add the dispatch entry which points to the function that implements the opcode.
 *
 * 4. Implement the actual opcode, nominally in the appropriate op_*.go file.
 */

// Constant describing instruction opcodes.
type Opcode int

const (
	Stop Opcode = iota // Stop must be the zero-th item.
	Add
	And
	BitAnd
	BitOr
	BitShift
	Call
	Coerce
	Copy
	Div
	Drop
	DropToMarker
	Dup
	Equal
	Exp
	Explode
	Flatten
	GreaterThan
	GreaterThanOrEqual
	LessThan
	LessThanOrEqual
	Load
	LoadIndex
	LoadSlice
	LoadThis
	Modulo
	Mul
	Negate
	NoOperation
	NotEqual
	Or
	Push
	ReadStack
	RequiredType
	SetThis
	StaticTyping
	Sub
	Swap

	// Everything from here on is a branch instruction, whose
	// operand must be present and is an integer instruction
	// address in the bytecode array. These instructions are
	// patched with offsets when code is appended.
	//
	// The first one in this list MIUST be BranchInstructions,
	// as it marks the start of the branch instructions, which
	// are instructions that can reference a bytecode address
	// as the operand.
	BranchInstructions
	Branch
	BranchTrue
	BranchFalse
)

var opcodeNames = map[Opcode]string{
	Add:                "Add",
	And:                "And",
	BitAnd:             "BitAnd",
	BitOr:              "BitOr",
	BitShift:           "BitShift",
	Branch:             "Branch",
	BranchFalse:        "BranchFalse",
	BranchTrue:         "BranchTrue",
	Call:               "Call",
	Coerce:             "Coerce",
	Copy:               "Copy",
	Div:                "Div",
	Drop:               "Drop",
	DropToMarker:       "DropToMarker",
	Dup:                "Dup",
	Equal:              "Equal",
	Exp:                "Exp",
	Explode:            "Explode",
	Flatten:            "Flatten",
	GreaterThan:        "GT",
	GreaterThanOrEqual: "GTEQ",
	LessThan:           "LT",
	LessThanOrEqual:    "LTEQ",
	Load:               "Load",
	LoadIndex:          "LoadIndex",
	LoadSlice:          "LoadSlice",
	LoadThis:           "LoadThis",
	Modulo:             "Modulo",
	Mul:                "Mul",
	Negate:             "Negate",
	NoOperation:        "NoOperation",
	NotEqual:           "NotEqual",
	Or:                 "Or",
	Push:               "Push",
	ReadStack:          "ReadStack",
	RequiredType:       "RequiredType",
	StaticTyping:       "StaticTyping",
	Stop:               "Stop",
	Sub:                "Sub",
	Swap:               "Swap",
}

func initializeDispatch() {
	if dispatch == nil {
		dispatch = dispatchMap{
			Add:                addByteCode,
			And:                andByteCode,
			BitAnd:             bitAndByteCode,
			BitOr:              bitOrByteCode,
			BitShift:           bitShiftByteCode,
			Branch:             branchByteCode,
			BranchFalse:        branchFalseByteCode,
			BranchTrue:         branchTrueByteCode,
			Call:               callByteCode,
			Coerce:             coerceByteCode,
			Div:                divideByteCode,
			Drop:               dropByteCode,
			DropToMarker:       dropToMarkerByteCode,
			Dup:                dupByteCode,
			Equal:              equalByteCode,
			Exp:                exponentByteCode,
			Explode:            explodeByteCode,
			GreaterThan:        greaterThanByteCode,
			GreaterThanOrEqual: greaterThanOrEqualByteCode,
			LessThan:           lessThanByteCode,
			LessThanOrEqual:    lessThanOrEqualByteCode,
			Load:               loadByteCode,
			LoadThis:           loadThisByteCode,
			Modulo:             moduloByteCode,
			Mul:                multiplyByteCode,
			Negate:             negateByteCode,
			NoOperation:        nil,
			NotEqual:           notEqualByteCode,
			Or:                 orByteCode,
			Push:               pushByteCode,
			ReadStack:          readStackByteCode,
			RequiredType:       requiredTypeByteCode,
			SetThis:            setThisByteCode,
			StaticTyping:       staticTypingByteCode,
			Stop:               stopByteCode,
			Sub:                subtractByteCode,
			Swap:               swapByteCode,
		}
	}
}
