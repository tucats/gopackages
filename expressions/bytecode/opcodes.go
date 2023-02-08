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
	AddressOf
	And
	Array
	BitAnd
	BitOr
	BitShift
	Call
	Coerce
	Copy
	CreateAndStore
	DeRef
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
	MakeArray
	MakeMap
	Modulo
	Mul
	Negate
	NoOperation
	NotEqual
	Or
	Push
	ReadStack
	RequiredType
	Return
	SetThis
	StaticTyping
	Store
	StoreAlways
	StoreGlobal
	StoreIndex
	StoreInto
	StoreViaPointer
	Struct
	Sub
	Swap
	SymbolCreate
	SymbolDelete
	SymbolOptCreate

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
	LocalCall
	RangeNext
	Try
)

var opcodeNames = map[Opcode]string{
	Add:                "Add",
	AddressOf:          "AddressOf",
	And:                "And",
	Array:              "Array",
	BitAnd:             "BitAnd",
	BitOr:              "BitOr",
	BitShift:           "BitShift",
	Branch:             "Branch",
	BranchFalse:        "BranchFalse",
	BranchTrue:         "BranchTrue",
	Call:               "Call",
	Coerce:             "Coerce",
	Copy:               "Copy",
	CreateAndStore:     "CreateAndStore",
	DeRef:              "DeRef",
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
	LocalCall:          "LocalCall",
	MakeArray:          "MakeArray",
	MakeMap:            "MakeMap",
	Modulo:             "Modulo",
	Mul:                "Mul",
	Negate:             "Negate",
	NoOperation:        "NoOperation",
	NotEqual:           "NotEqual",
	Or:                 "Or",
	Push:               "Push",
	ReadStack:          "ReadStack",
	RequiredType:       "RequiredType",
	Return:             "Return",
	StaticTyping:       "StaticTyping",
	Stop:               "Stop",
	Store:              "Store",
	StoreAlways:        "StoreAlways",
	StoreGlobal:        "StoreGlobal",
	StoreIndex:         "StoreIndex",
	StoreInto:          "StoreInto",
	StoreViaPointer:    "StorePointer",
	Struct:             "Struct",
	Sub:                "Sub",
	Swap:               "Swap",
	SymbolCreate:       "SymbolCreate",
	SymbolDelete:       "SymbolDelete",
	SymbolOptCreate:    "SymbolOptCreate",
	Try:                "Try",
}

func initializeDispatch() {
	if dispatch == nil {
		dispatch = dispatchMap{
			Add:                addByteCode,
			AddressOf:          addressOfByteCode,
			And:                andByteCode,
			Array:              arrayByteCode,
			BitAnd:             bitAndByteCode,
			BitOr:              bitOrByteCode,
			BitShift:           bitShiftByteCode,
			Branch:             branchByteCode,
			BranchFalse:        branchFalseByteCode,
			BranchTrue:         branchTrueByteCode,
			Call:               callByteCode,
			Coerce:             coerceByteCode,
			CreateAndStore:     createAndStoreByteCode,
			DeRef:              deRefByteCode,
			Div:                divideByteCode,
			Drop:               dropByteCode,
			DropToMarker:       dropToMarkerByteCode,
			Dup:                dupByteCode,
			Equal:              equalByteCode,
			Exp:                exponentByteCode,
			Explode:            explodeByteCode,
			Flatten:            flattenByteCode,
			GreaterThan:        greaterThanByteCode,
			GreaterThanOrEqual: greaterThanOrEqualByteCode,
			LessThan:           lessThanByteCode,
			LessThanOrEqual:    lessThanOrEqualByteCode,
			Load:               loadByteCode,
			LoadIndex:          loadIndexByteCode,
			LoadSlice:          loadSliceByteCode,
			LoadThis:           loadThisByteCode,
			MakeArray:          makeArrayByteCode,
			MakeMap:            makeMapByteCode,
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
			Store:              storeByteCode,
			StoreAlways:        storeAlwaysByteCode,
			StoreGlobal:        storeGlobalByteCode,
			StoreIndex:         storeIndexByteCode,
			StoreInto:          storeIntoByteCode,
			StoreViaPointer:    storeViaPointerByteCode,
			Struct:             structByteCode,
			Sub:                subtractByteCode,
			Swap:               swapByteCode,
			SymbolCreate:       symbolCreateByteCode,
			SymbolDelete:       symbolDeleteByteCode,
			SymbolOptCreate:    symbolCreateIfByteCode,
		}
	}
}
