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

	Add
	And
	ArgCheck
	Array
	Auth
	Call
	ClassMember
	Coerce
	Constant
	Copy
	Div
	Drop
	DropToMarker
	Dup
	Equal
	Exp
	Flatten
	GreaterThan
	GreaterThanOrEqual
	LessThan
	LessThanOrEqual
	Load
	LoadIndex
	LoadSlice
	Log
	MakeArray
	Member
	Mul
	Negate
	Newline
	NotEqual
	Or
	Panic
	PopScope
	Print
	Push
	PushScope
	Response
	Return
	Say
	StackCheck
	StaticTyping
	Store
	StoreAlways
	StoreGlobal
	StoreIndex
	Struct
	Sub
	Swap
	SymbolCreate
	SymbolDelete
	SymbolOptCreate
	Template
	This
	Try
	TryPop

	// Everything from here on is a branch instruction, whose
	// operand must be present and is an integer instruction
	// address in the bytecode array
	BranchInstructions = iota + BranchInstruction
	Branch
	BranchTrue
	BranchFalse
	LocalCall

	// After this value, additional user branch instructions are
	// can be defined.
	UserBranchInstructions
)

var opcodeNames = map[int]string{
	Add:                "Add",
	And:                "And",
	ArgCheck:           "ArgCheck",
	Array:              "Array",
	AtLine:             "AtLine",
	Auth:               "Auth",
	Branch:             "Branch",
	BranchFalse:        "BranchFalse",
	BranchTrue:         "BranchTrue",
	Call:               "Call",
	ClassMember:        "ClassMember",
	Coerce:             "Coerce",
	Constant:           "Constant",
	Copy:               "Copy",
	Div:                "Div",
	Drop:               "Drop",
	DropToMarker:       "DropToMarker",
	Dup:                "Dup",
	Equal:              "Equal",
	Exp:                "Exp",
	Flatten:            "Flatten",
	GreaterThan:        "GreaterThan",
	GreaterThanOrEqual: "GreaterThanOrEqual",
	LessThan:           "LessThan",
	LessThanOrEqual:    "LessThanOrEqual",
	Load:               "Load",
	LoadIndex:          "LoadIndex",
	LoadSlice:          "LoadSlice",
	LocalCall:          "LocalCall",
	Log:                "Log",
	MakeArray:          "MakeArray",
	Member:             "Member",
	Mul:                "Mul",
	Negate:             "Negate",
	Newline:            "Newline",
	NotEqual:           "NotEqual",
	Or:                 "Or",
	Panic:              "Panic",
	PopScope:           "PopScope",
	Print:              "Print",
	Push:               "Push",
	PushScope:          "PushScope",
	Response:           "Response",
	Return:             "Return",
	Say:                "Say",
	StackCheck:         "StackCheck",
	StaticTyping:       "StaticTyping",
	Stop:               "Stop",
	Store:              "Store",
	StoreAlways:        "StoreAlways",
	StoreGlobal:        "StoreGlobal",
	StoreIndex:         "StoreIndex",
	Struct:             "Struct",
	Sub:                "Sub",
	Swap:               "Swap",
	SymbolCreate:       "SymbolCreate",
	SymbolDelete:       "SymbolDelete",
	SymbolOptCreate:    "SymbolOptCreate",
	Template:           "Template",
	This:               "This",
	Try:                "Try",
	TryPop:             "TryPop",
}

func initializeDispatch() {
	if dispatch == nil {
		dispatch = DispatchMap{
			Add:                AddOpcode,
			And:                AndOpcode,
			ArgCheck:           ArgCheckOpcode,
			Array:              ArrayOpcode,
			AtLine:             AtLineOpcode,
			Auth:               AuthOpcode,
			Branch:             BranchOpcode,
			BranchFalse:        BranchFalseOpcode,
			BranchTrue:         BranchTrueOpcode,
			Call:               CallOpcode,
			ClassMember:        ClassMemberOpcode,
			Coerce:             CoerceOpcode,
			Constant:           ConstantOpcode,
			Copy:               CopyOpcode,
			Div:                DivOpcode,
			Drop:               DropOpcode,
			DropToMarker:       DropToMarkerOpcode,
			Dup:                DupOpcode,
			Equal:              EqualOpcode,
			Exp:                ExpOpcode,
			Flatten:            FlattenOpcode,
			GreaterThan:        GreaterThanOpcode,
			GreaterThanOrEqual: GreaterThanOrEqualOpcode,
			LessThan:           LessThanOpcode,
			LessThanOrEqual:    LessThanOrEqualOpcode,
			Load:               LoadOpcode,
			LoadIndex:          LoadIndexOpcode,
			LoadSlice:          LoadSliceOpcode,
			LocalCall:          LocalCallOpcode,
			Log:                LogOpcode,
			MakeArray:          MakeArrayOpcode,
			Member:             MemberOpcode,
			Mul:                MulOpcode,
			Negate:             NegateOpcode,
			Newline:            NewlineOpcode,
			NotEqual:           NotEqualOpcode,
			Or:                 OrOpcode,
			Panic:              PanicOpcode,
			PopScope:           PopScopeOpcode,
			Print:              PrintOpcode,
			Push:               PushOpcode,
			PushScope:          PushScopeOpcode,
			Response:           ResponseOpcode,
			Return:             ReturnOpcode,
			Say:                SayOpcode,
			StackCheck:         StackCheckOpcode,
			StaticTyping:       StaticTypingOpcode,
			Stop:               StopOpcode,
			Store:              StoreOpcode,
			StoreAlways:        StoreAlwaysOpcode,
			StoreGlobal:        StoreGlobalOpcode,
			StoreIndex:         StoreIndexOpcode,
			Struct:             StructOpcode,
			Sub:                SubOpcode,
			Swap:               SwapOpcode,
			SymbolCreate:       SymbolCreateOpcode,
			SymbolDelete:       SymbolDeleteOpcode,
			SymbolOptCreate:    SymbolOptCreateOpcode,
			Template:           TemplateOpcode,
			This:               ThisOpcode,
			Try:                TryOpcode,
			TryPop:             TryPopOpcode,
		}
	}
}
