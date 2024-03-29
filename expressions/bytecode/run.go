package bytecode

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/tucats/gopackages/app-cli/ui"
	"github.com/tucats/gopackages/defs"
	"github.com/tucats/gopackages/errors"
	"github.com/tucats/gopackages/expressions/symbols"
	"github.com/tucats/gopackages/i18n"
)

// opcodeHandler defines a function that implements an opcode.
type opcodeHandler func(b *Context, i interface{}) error

// dispatchMap is a map that is used to locate the function for an opcode.
type dispatchMap map[Opcode]opcodeHandler

var dispatch dispatchMap
var dispatchMux sync.Mutex
var waitGroup sync.WaitGroup

// growStackBy indicates the number of elements to add to the stack when
// it runs out of space.
const growStackBy = 50

func (c *Context) GetName() string {
	if c.bc != nil {
		return c.bc.name
	}

	return defs.Main
}

func (c *Context) StepOver(b bool) {
	c.stepOver = b
}

func (c *Context) GetSymbols() *symbols.SymbolTable {
	return c.symbols
}

// Run executes a bytecode context.
func (c *Context) Run() error {
	return c.RunFromAddress(0)
}

// Used to resume execution after an event like the debugger being invoked.
func (c *Context) Resume() error {
	return c.RunFromAddress(c.programCounter)
}

func (c *Context) IsRunning() bool {
	return c.running
}

// RunFromAddress executes a bytecode context from a given starting address.
func (c *Context) RunFromAddress(addr int) error {
	var err error

	// Make sure globals are initialized. Because this updates a global, let's
	// do it in a thread-safe fashion.
	dispatchMux.Lock()
	initializeDispatch()
	dispatchMux.Unlock()

	// Reset the runtime context.
	c.programCounter = addr
	c.running = true

	ui.Log(ui.TraceLogger, "*** Tracing %s (%d)  ", c.name, c.threadID)

	// Loop over the bytecodes and run.
	for c.running {
		if c.programCounter >= len(c.bc.instructions) {
			c.running = false

			break
		}

		i := c.bc.instructions[c.programCounter]

		atomic.AddInt64(&InstructionsExecuted, 1)

		if c.Tracing() {
			instruction := FormatInstruction(i)

			stack := c.formatStack(c.symbols, c.fullStackTrace)
			if !c.fullStackTrace && len(stack) > 80 {
				stack = stack[:80]
			}

			if len(instruction) > 30 {
				ui.Log(ui.TraceLogger, "(%d) %18s %3d: %s",
					c.threadID, c.GetModuleName(), c.programCounter, instruction)
				ui.Log(ui.TraceLogger, "(%d) %18s %3s  %-30s stack[%2d]: %s",
					c.threadID, " ", " ", " ", c.stackPointer, stack)
			} else {
				ui.Log(ui.TraceLogger, "(%d) %18s %3d: %-30s stack[%2d]: %s",
					c.threadID, c.GetModuleName(), c.programCounter, instruction, c.stackPointer, stack)
			}
		}

		c.programCounter = c.programCounter + 1

		imp, found := dispatch[i.Operation]
		if !found {
			return c.error(errors.ErrUnimplementedInstruction).Context(i.Operation)
		}

		err = imp(c, i.Operand)
		if err != nil {
			if !errors.Equals(err, errors.ErrSignalDebugger) && !errors.Equals(err, errors.ErrStop) {
				ui.Log(ui.TraceLogger, "(%d)  *** Return error: %s", c.threadID, err)
			}

			if err != nil {
				err = errors.NewError(err)
			}

			return err
		}
	}

	ui.Log(ui.TraceLogger, "*** End tracing %s (%d) ", c.name, c.threadID)

	if err != nil {
		return errors.NewError(err)
	}

	return nil
}

// GoRoutine allows calling a named function as a go routine, using arguments. The invocation
// of GoRoutine should be in a "go" statement to run the code.
func GoRoutine(fName string, parentCtx *Context, args []interface{}) {
	parentCtx.mux.RLock()
	parentSymbols := parentCtx.symbols
	parentCtx.mux.RUnlock()

	err := parentCtx.error(errors.ErrInvalidFunctionCall)

	ui.Log(ui.TraceLogger, "--> Starting Go routine \"%s\"", fName)
	ui.Log(ui.TraceLogger, "--> Argument list: %#v", args)

	// Locate the bytecode for the function. It must be a symbol defined as bytecode.
	if fCode, ok := parentSymbols.Get(fName); ok {
		if bc, ok := fCode.(*ByteCode); ok {
			bc.Disasm()
			// Create a new stream whose job is to invoke the function by name.
			callCode := New("go " + fName)
			callCode.Emit(Load, fName)

			for _, arg := range args {
				callCode.Emit(Push, arg)
			}

			callCode.Emit(Call, len(args))

			// Make a new table that is parently only to the root table (for access to
			// packages). Copy the function definition into this new table so the invocation
			// of the function within the native go routine can locate it.
			functionSymbols := symbols.NewChildSymbolTable("Go routine "+fName, parentSymbols.SharedParent())
			functionSymbols.SetAlways(fName, bc)

			ctx := NewContext(functionSymbols, callCode)
			err = parentCtx.error(ctx.Run())

			waitGroup.Done()
		}
	}

	if err != nil && !err.Is(errors.ErrStop) {
		fmt.Printf("%s\n", i18n.E("go.error", map[string]interface{}{"name": fName, "err": err}))

		ui.Log(ui.TraceLogger, "--> Go routine invocation ends with %v", err)
		os.Exit(55)
	}
}
