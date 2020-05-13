package bytecode

import (
	"errors"
	"strconv"
)

// OpcodeHandler defines a function that implements an opcode
type OpcodeHandler func(b *Context, i *I) error

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
func (c *Context) Run() error {

	var err error

	c.pc = 0
	c.running = true
	c.sp = 0

	for c.running {

		if c.pc > len(c.bc.opcodes) {
			return errors.New("ran off end of incomplete bytecode")
		}
		i := c.bc.opcodes[c.pc]
		c.pc = c.pc + 1

		imp, found := dispatch[i.Opcode]
		if !found {
			return errors.New("inimplemented instruction: " + strconv.Itoa(i.Opcode))
		}
		err = imp(c, &i)
		if err != nil {
			return err
		}
	}

	return err
}
