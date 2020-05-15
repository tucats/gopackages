package bytecode

import (
	"errors"
	"fmt"

	"github.com/tucats/gopackages/util"
)

// GrowOpcodesBy indicates the number of elements to add to the
// opcode array when storage is exhausted in the current array.
const GrowOpcodesBy = 50

// InitialOpcodeSize is the initial size of the emit buffer
const InitialOpcodeSize = 20

// InitialStackSize is the initial stack size.
const InitialStackSize = 100

// BranchInstruction is the mininum value for a branch instruction, which has
// special meaning during relocation and linking
const BranchInstruction = 65536

// BuiltinInstructions defines the lowest number (other than Stop) of the
// builtin instructions provided by the bytecode package. User instructions
// are added between 1 and this value.
const BuiltinInstructions = BranchInstruction - 2048

// Constant describing instruction opcodes
const (
	Stop = 0
	Call = iota + BuiltinInstructions
	Push
	Drop
	Add
	Sub
	Div
	Mul
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
	LoadIndex
	StoreIndex
	Struct
	Member
	Print
	Newline

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

// I contains the information about a single bytecode instruction.
type I struct {
	Opcode  int
	Operand interface{}
}

// ByteCode contains the context of the execution of a bytecode stream.
type ByteCode struct {
	Name    string
	opcodes []I
	emitPos int
	Symbols SymbolTable
}

// New generates and initializes a new bytecode
func New(name string) *ByteCode {

	bc := ByteCode{
		Name:    name,
		opcodes: make([]I, InitialOpcodeSize),
		emitPos: 0,
		Symbols: SymbolTable{},
	}

	return &bc
}

// Emit0 emits an instruction with no operand
func (b *ByteCode) Emit0(opcode int) {
	b.Emit(opcode, nil)
}

// Emit emits a single instruction
func (b *ByteCode) Emit(opcode int, operand interface{}) {
	if b.emitPos >= len(b.opcodes) {
		b.opcodes = append(b.opcodes, make([]I, GrowOpcodesBy)...)
	}
	i := I{Opcode: opcode, Operand: operand}
	b.opcodes[b.emitPos] = i
	b.emitPos = b.emitPos + 1
}

// Mark returns the address of the instruction about to be emitted.
func (b *ByteCode) Mark() int {
	return b.emitPos
}

// SetAddressHere sets the current address as the target of the marked
// instruction
func (b *ByteCode) SetAddressHere(mark int) error {
	return b.SetAddress(mark, b.emitPos)
}

// SetAddress sets the given value as the target of the marked
// instruction
func (b *ByteCode) SetAddress(mark int, address int) error {

	if mark > b.emitPos || mark < 0 {
		return errors.New("invalid marked position")
	}
	i := b.opcodes[mark]
	i.Operand = address
	b.opcodes[mark] = i
	return nil
}

// Append appends another bytecode to the current bytecode,
// and updates all the link references.
func (b *ByteCode) Append(a *ByteCode) {

	base := b.emitPos

	for _, i := range a.opcodes[:a.emitPos] {
		if i.Opcode > BranchInstructions {
			i.Operand = util.GetInt(i.Operand) + base
		}
		b.Emit(i.Opcode, i.Operand)
	}
}

// DefineInstruction adds a user-defined instruction to the bytecode
// set.
func DefineInstruction(opcode int, name string, implementation OpcodeHandler) error {

	// First, make sure this isn't a duplicate
	if _, found := dispatch[opcode]; found {
		return fmt.Errorf("opcode already defined: %d", opcode)
	}

	opcodeNames[opcode] = name
	dispatch[opcode] = implementation

	return nil
}

// Run generates a one-time context for executing this bytecode.
func (b *ByteCode) Run(s *SymbolTable) error {
	c := NewContext(s, b)
	return c.Run()
}

// Call generates a one-time context for executing this bytecode,
// and returns a value as well as an error.
func (b *ByteCode) Call(s *SymbolTable) (interface{}, error) {
	c := NewContext(s, b)
	err := c.Run()
	if err != nil {
		return nil, err
	}

	return c.Pop()
}

// Opcodes returns the opcode list for this byteocde array
func (b *ByteCode) Opcodes() []I {
	return b.opcodes[:b.emitPos]
}
