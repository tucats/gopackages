package bytecode

import (
	"strconv"
	"sync"

	"github.com/tucats/gopackages/app-cli/ui"
)

// OpcodeHandler defines a function that implements an opcode
type OpcodeHandler func(b *Context, i interface{}) error

// DispatchMap is a map that is used to locate the function for an opcode
type DispatchMap map[int]OpcodeHandler

var dispatch DispatchMap
var dispatchMux sync.Mutex

// GrowStackBy indicates the number of eleemnts to add to the stack when
// it runs out of space.
const GrowStackBy = 50

func initializeDispatch() {
	if dispatch == nil {
		dispatch = DispatchMap{
			Stop:               StopOpcode,
			AtLine:             AtLineOpcode,
			Push:               PushOpcode,
			Array:              ArrayOpcode,
			LoadIndex:          LoadIndexOpcode,
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

	// Make sure globals are initialized. Becuase this updates a global, let's
	// do it in a thread-safe fashion.
	dispatchMux.Lock()
	initializeDispatch()
	dispatchMux.Unlock()

	// Reset the runtime context
	c.pc = addr
	c.running = true
	c.sp = 0

	// Make sure the opcode array ends in a Stop operation
	if c.bc.emitPos == 0 || c.bc.opcodes[c.bc.emitPos-1].Opcode != Stop {
		ui.Debug("Adding trailing Stop opcode")
		c.bc.Emit1(Stop)
	}

	if c.Tracing {
		ui.Debug("*** Tracing " + c.Name)
	}
	// Loop over the bytecodes and run.
	for c.running {

		if c.pc > len(c.bc.opcodes) {
			return c.NewError("ran off end of incomplete bytecode")
		}

		i := c.bc.opcodes[c.pc]
		if c.Tracing {
			s := FormatInstruction(i)
			s2 := FormatStack(c.stack[:c.sp])
			if len(s2) > 50 {
				s2 = s2[:50]
			}
			ui.Debug("%5d: %-30s stack: %s", c.pc, s, s2)
		}
		c.pc = c.pc + 1

		imp, found := dispatch[i.Opcode]
		if !found {
			return c.NewStringError("unimplemented instruction", strconv.Itoa(i.Opcode))
		}
		err = imp(c, i.Operand)
		if err != nil {

			text := err.Error()

			// See if we are in a try/catch block. IF there is a Try/Catch stack
			// and the jump point on top is non-zero, then we can transfer control.
			if len(c.try) > 0 && c.try[len(c.try)-1] > 0 {
				c.pc = c.try[len(c.try)-1]

				// Zero out the jump point for this try/catch block so recursive
				// errors don't occur.
				c.try[len(c.try)-1] = 0
				c.symbols.SetAlways("_error", text)
				if c.Tracing {
					ui.Debug("*** Branch to %d on error: %s", c.pc, text)
				}
			} else {
				if c.Tracing {
					ui.Debug("*** Return error: %s", text)
				}
				return err
			}
		}
	}
	if c.Tracing {
		ui.Debug("*** End tracing " + c.Name)
	}

	return err
}
