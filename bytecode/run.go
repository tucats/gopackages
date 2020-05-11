package bytecode

import (
	"errors"
	"strconv"
)

// OpcodeHandler defines a function that implements an opcode
type OpcodeHandler func(b *ByteCode, i *I) error

// DispatchMap is a map that is used to locate the function for an opcode
type DispatchMap map[int]OpcodeHandler

var dispatch = DispatchMap{
	Stop:               StopOpcode,
	Push:               PushOpcode,
	Array:              ArrayOpcode,
	Index:              IndexOpcode,
	Add:                AddOpcode,
	Sub:                SubOpcode,
	Mul:                MulOpcode,
	Div:                DivOpcode,
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
}

// GrowStackBy indicates the number of eleemnts to add to the stack when
// it runs out of space.
const GrowStackBy = 50

// Run executes a bytecode context
func (b *ByteCode) Run(symbols map[string]interface{}) error {

	var err error

	b.pc = 0
	b.running = true
	b.sp = 0
	b.symbols = symbols

	for b.running {

		if b.pc > len(b.opcodes) {
			return errors.New("ran off end of incomplete bytecode")
		}
		i := b.opcodes[b.pc]
		b.pc = b.pc + 1

		imp, found := dispatch[i.Opcode]
		if !found {
			return errors.New("inimplemented instruction: " + strconv.Itoa(i.Opcode))
		}
		err = imp(b, &i)
		if err != nil {
			return err
		}
	}

	return err
}

// Pop removes the top-most item from the stack
func (b *ByteCode) Pop() (interface{}, error) {
	if b.sp <= 0 || len(b.stack) < b.sp {
		return nil, errors.New("stack underflow")
	}

	b.sp = b.sp - 1
	v := b.stack[b.sp]
	return v, nil
}

// Push puts a new items on the stack
func (b *ByteCode) Push(v interface{}) error {

	if b.sp >= len(b.stack) {
		b.stack = append(b.stack, make([]interface{}, GrowStackBy)...)
	}
	b.stack[b.sp] = v
	b.sp = b.sp + 1
	return nil
}
