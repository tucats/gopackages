package bytecode

import "errors"

// GrowOpcodesBy indicates the number of elements to add to the
// opcode array when storage is exhausted in the current array.
const GrowOpcodesBy = 50

// InitialOpcodeSize is the initial size of the emit buffer
const InitialOpcodeSize = 20

// InitialStackSize is the initial stack size.
const InitialStackSize = 100

// Constant describing instruction opcodes
const (
	Stop = iota
	Call
	Push
	Array
	Add
	Sub
	Div
	Mul
	And
	Or
	Negate
	Branch
	BranchTrue
	BranchFalse
	Equal
	NotEqual
	GreaterThan
	LessThan
	GreaterThanOrEqual
	LessThanOrEqual
	Load
	Store
	Index
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
}

// New generates and initializes a new bytecode
func New(name string) *ByteCode {

	bc := ByteCode{
		Name:    name,
		opcodes: make([]I, InitialOpcodeSize),
		emitPos: 0,
	}

	return &bc
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
