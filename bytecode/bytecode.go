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
	Add
	Sub
	Div
	Mul
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
)

// I contains the information about a single bytecode instruction.
type I struct {
	opcode  int
	operand interface{}
}

// ByteCode contains the context of the execution of a bytecode stream.
type ByteCode struct {
	Name    string
	opcodes []I
	emitPos int
	pc      int
	stack   []interface{}
	sp      int
	running bool
	symbols map[string]interface{}
}

// New generates and initializes a new bytecode
func New(name string) *ByteCode {

	bc := ByteCode{
		Name:    name,
		opcodes: make([]I, InitialOpcodeSize),
		stack:   make([]interface{}, InitialStackSize),
		emitPos: 0,
		running: false,
		pc:      0,
		sp:      0,
	}

	return &bc
}

// Emit emits a single instruction
func (b *ByteCode) Emit(i I) {
	if b.emitPos >= len(b.opcodes) {
		b.opcodes = append(b.opcodes, make([]I, GrowOpcodesBy)...)
	}
	b.opcodes[b.emitPos] = i
	b.emitPos = b.emitPos + 1
}

// Mark returns the address of the instruction about to be emitted.
func (b *ByteCode) Mark() int {
	return b.emitPos
}

// SetAddress sets the given value as the target of the marked
// instruction
func (b *ByteCode) SetAddress(mark int, address int) error {
	if mark > b.emitPos || mark < 0 {
		return errors.New("invalid marked position")
	}
	i := b.opcodes[mark]
	i.operand = address
	b.opcodes[mark] = i
	return nil
}

// Get retrieves a symbol value from the symbol table
func (b *ByteCode) Get(name string) interface{} {
	return b.symbols[name]
}

// Set sets a symbol value in the symbol table
func (b *ByteCode) Set(name string, value interface{}) {
	b.symbols[name] = value
}
