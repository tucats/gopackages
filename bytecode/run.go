package bytecode

import (
	"errors"
	"strconv"

	"github.com/tucats/gopackages/app-cli/ui"
)

// OpcodeHandler defines a function that implements an opcode
type OpcodeHandler func(b *Context, i *I) error

// DispatchMap is a map that is used to locate the function for an opcode
type DispatchMap map[int]OpcodeHandler

var dispatch DispatchMap

// GrowStackBy indicates the number of eleemnts to add to the stack when
// it runs out of space.
const GrowStackBy = 50

func initializeDispatch() {
	if dispatch == nil {
		dispatch = DispatchMap{
			Stop:               StopOpcode,
			Push:               PushOpcode,
			Array:              ArrayOpcode,
			Index:              IndexOpcode,
			Struct:             StructOpcode,
			Member:             MemberOpcode,
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
			Print:              PrintOpcode,
			Newline:            NewlineOpcode,
		}
	}
}

// Run executes a bytecode context
func (c *Context) Run() error {
	return c.RunFromAddress(0)
}

// RunFromAddress executes a bytecode context from a given starting address.
func (c *Context) RunFromAddress(addr int) error {

	var err error

	// Make sure globals are initialized
	initializeDispatch()

	// Reset the runtime context
	c.pc = addr
	c.running = true
	c.sp = 0

	// Make sure the opcode array ends in a Stop operation
	if c.bc.emitPos == 0 || c.bc.opcodes[c.bc.emitPos-1].Opcode != Stop {
		ui.Debug("Adding trailing Stop opcode")
		c.bc.Emit(Stop, nil)
	}

	// Loop over the bytecodes and run.
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
