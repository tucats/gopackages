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

// Constant describing instruction opcodes
const (
	Stop   = 0
	AtLine = iota + BuiltinInstructions
	Call
	ArgCheck
	Push
	Drop
	Dup
	Copy
	Swap
	Add
	Sub
	Div
	Mul
	Exp
	And
	Or
	Negate
	Equal
	NotEqual
	GreaterThan
	LessThan
	GreaterThanOrEqual
	LessThanOrEqual
	Load
	Store
	Array
	MakeArray
	LoadIndex
	StoreIndex
	LoadSlice
	Struct
	Member
	ClassMember
	Print
	Say
	Newline
	SymbolDelete
	SymbolCreate
	Constant
	PushScope
	PopScope
	Try
	TryPop
	Coerce
	This
	Panic
	Template
	Return

	// Everything from here on is a branch instruction, whose
	// operand must be present and is an integer instruction
	// address in the bytecode array
	BranchInstructions = iota + BranchInstruction
	Branch
	BranchTrue
	BranchFalse

	// After this value, additional user branch instructions are
	// can be defined.
	UserBranchInstructions
)

var opcodeNames = map[int]string{
	Stop:               "Stop",
	Call:               "Call",
	Push:               "Push",
	Add:                "Add",
	Sub:                "Sub",
	Div:                "Div",
	Mul:                "Mul",
	Exp:                "Exp",
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
	LoadSlice:          "LoadSlice",
	StoreIndex:         "StoreIndex",
	Struct:             "Struct",
	Member:             "Member",
	Print:              "Print",
	Say:                "Say",
	Newline:            "Newline",
	Drop:               "Drop",
	AtLine:             "AtLine",
	MakeArray:          "MakeArray",
	SymbolDelete:       "SymbolDelete",
	SymbolCreate:       "SymbolCreate",
	PushScope:          "PushScope",
	PopScope:           "PopScope",
	Constant:           "Constant",
	Try:                "Try",
	TryPop:             "TryPop",
	Coerce:             "Coerce",
	ArgCheck:           "ArgCheck",
	This:               "This",
	Dup:                "Dup",
	Copy:               "Copy",
	Swap:               "Swap",
	ClassMember:        "ClassMember",
	Panic:              "Panic",
	Template:           "Template",
	Return:             "Return",
}

func initializeDispatch() {
	if dispatch == nil {
		dispatch = DispatchMap{
			Panic:              PanicOpcode,
			Stop:               StopOpcode,
			AtLine:             AtLineOpcode,
			Push:               PushOpcode,
			Array:              ArrayOpcode,
			LoadIndex:          LoadIndexOpcode,
			LoadSlice:          LoadSliceOpcode,
			StoreIndex:         StoreIndexOpcode,
			Struct:             StructOpcode,
			Member:             MemberOpcode,
			Add:                AddOpcode,
			Sub:                SubOpcode,
			Mul:                MulOpcode,
			Div:                DivOpcode,
			Exp:                ExpOpcode,
			And:                AndOpcode,
			Or:                 OrOpcode,
			Negate:             NegateOpcode,
			Call:               CallOpcode,
			Load:               LoadOpcode,
			Store:              StoreOpcode,
			Branch:             BranchOpcode,
			BranchTrue:         BranchTrueOpcode,
			BranchFalse:        BranchFalseOpcode,
			Equal:              EqualOpcode,
			NotEqual:           NotEqualOpcode,
			LessThan:           LessThanOpcode,
			LessThanOrEqual:    LessThanOrEqualOpcode,
			GreaterThan:        GreaterThanOpcode,
			GreaterThanOrEqual: GreaterThanOrEqualOpcode,
			Print:              PrintOpcode,
			Say:                SayOpcode,
			Newline:            NewlineOpcode,
			Drop:               DropOpcode,
			MakeArray:          MakeArrayOpcode,
			SymbolDelete:       SymbolDeleteOpcode,
			SymbolCreate:       SymbolCreateOpcode,
			PushScope:          PushScopeOpcode,
			PopScope:           PopScopeOpcode,
			Constant:           ConstantOpcode,
			Try:                TryOpcode,
			TryPop:             TryPopOpcode,
			Coerce:             CoerceOpcode,
			ArgCheck:           ArgCheckOpcode,
			This:               ThisOpcode,
			Dup:                DupOpcode,
			Copy:               CopyOpcode,
			Swap:               SwapOpcode,
			ClassMember:        ClassMemberOpcode,
			Template:           TemplateOpcode,
			Return:             ReturnOpcode,
		}
	}
}
