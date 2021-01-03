package bytecode

import (
	"strconv"
	"sync"

	"github.com/tucats/gopackages/app-cli/ui"
	"github.com/tucats/gopackages/util"
)

// OpcodeHandler defines a function that implements an opcode
type OpcodeHandler func(b *Context, i interface{}) error

// DispatchMap is a map that is used to locate the function for an opcode
type DispatchMap map[Instruction]OpcodeHandler

var dispatch DispatchMap
var dispatchMux sync.Mutex

// GrowStackBy indicates the number of eleemnts to add to the stack when
// it runs out of space.
const GrowStackBy = 50

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

	// Make sure the opcode array ends in a Stop operation so we can never
	// shoot off the end of the bytecode.
	if c.bc.emitPos == 0 || c.bc.opcodes[c.bc.emitPos-1].Operation != Stop {
		c.bc.Emit(Stop)
	}

	if c.Tracing {
		ui.Debug(ui.ByteCodeLogger, "*** Tracing "+c.Name)
	}

	fullStackListing := util.GetBool(c.GetConfig("full_stack_listing"))

	// Loop over the bytecodes and run.
	for c.running {

		if c.pc > len(c.bc.opcodes) {
			break
		}

		i := c.bc.opcodes[c.pc]
		if c.Tracing {
			s := FormatInstruction(i)
			s2 := FormatStack(c.stack[:c.sp], fullStackListing)
			if !fullStackListing && len(s2) > 50 {
				s2 = s2[:50]
			}
			ui.Debug(ui.ByteCodeLogger, "%5d: %-30s stack[%2d]: %s", c.pc, s, c.sp, s2)
		}
		c.pc = c.pc + 1

		imp, found := dispatch[i.Operation]
		if !found {
			return c.NewError(UnimplementedInstructionError, strconv.Itoa(int(i.Operation)))
		}
		err = imp(c, i.Operand)
		if err != nil {

			text := err.Error()

			// See if we are in a try/catch block. IF there is a Try/Catch stack
			// and the jump point on top is non-zero, then we can transfer control.
			// Note that if the error was fatal, the running flag is turned off, which
			// prevents the try block from being honored (i.e. you cannot catch a fatal
			// error)
			if len(c.try) > 0 && c.try[len(c.try)-1] > 0 && c.running {
				c.pc = c.try[len(c.try)-1]

				// Zero out the jump point for this try/catch block so recursive
				// errors don't occur.
				c.try[len(c.try)-1] = 0
				_ = c.symbols.SetAlways("_error", text)
				if c.Tracing {
					ui.Debug(ui.ByteCodeLogger, "*** Branch to %d on error: %s", c.pc, text)
				}
			} else {
				if c.Tracing {
					ui.Debug(ui.ByteCodeLogger, "*** Return error: %s", text)
				}
				return err
			}
		}
	}
	if c.Tracing {
		ui.Debug(ui.ByteCodeLogger, "*** End tracing "+c.Name)
	}

	return err
}
