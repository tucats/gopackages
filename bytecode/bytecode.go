package bytecode

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
