package bytecode

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/tucats/gopackages/functions"
	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/util"
)

/******************************************\
*                                         *
*        F L O W   C O N T R O L          *
*                                         *
\******************************************/

// StopImpl instruction processor causes the current execution context to
// stop executing immediately.
func StopImpl(c *Context, i interface{}) error {
	c.running = false
	return nil
}

// PanicImpl instruction processor generates an error. The boolean flag is used
// to indicate if this is a fatal error that stops Ego, versus a user error.
func PanicImpl(c *Context, i interface{}) error {
	c.running = !util.GetBool(i)
	strValue, err := c.Pop()
	if err != nil {
		return err
	}
	msg := util.GetString(strValue)
	return c.NewError(msg)
}

// AtLineImpl instruction processor. This identifies the start of a new statement,
// and tags the line number from the source where this was found. This is used
// in error messaging, primarily.
func AtLineImpl(c *Context, i interface{}) error {

	c.line = util.GetInt(i)
	// If we are tracing, put that out now.
	if c.tokenizer != nil {
		fmt.Printf("%d:  %s\n", c.line, c.tokenizer.GetLine(c.line))
	}
	return nil
}

// BranchFalseImpl instruction processor branches to the instruction named in
// the operand if the top-of-stack item is a boolean FALSE value. Otherwise,
// execution continues with the next instruction.
func BranchFalseImpl(c *Context, i interface{}) error {

	// Get test value
	v, err := c.Pop()
	if err != nil {
		return err
	}

	// Get destination
	address := util.GetInt(i)
	if address < 0 || address > c.bc.emitPos {
		return c.NewError(InvalidBytecodeAddress)
	}

	if !util.GetBool(v) {
		c.pc = address
	}
	return nil
}

// BranchImpl instruction processor branches to the instruction named in
// the operand.
func BranchImpl(c *Context, i interface{}) error {

	// Get destination
	address := util.GetInt(i)
	if address < 0 || address > c.bc.emitPos {
		return c.NewError(InvalidBytecodeAddress)
	}

	c.pc = address
	return nil
}

// BranchTrueImpl instruction processor branches to the instruction named in
// the operand if the top-of-stack item is a boolean TRUE value. Otherwise,
// execution continues with the next instruction.
func BranchTrueImpl(c *Context, i interface{}) error {

	// Get test value
	v, err := c.Pop()
	if err != nil {
		return err
	}

	// Get destination
	address := util.GetInt(i)
	if address < 0 || address > c.bc.emitPos {
		return c.NewError(InvalidBytecodeAddress)
	}

	if util.GetBool(v) {
		c.pc = address
	}
	return nil
}

// LocalCallImpl runs a subroutine (a function that has no parameters and
// no return value) that is compiled into the same bytecode as the current
// instruction stream. This is used to implement defer statement blocks, for
// example, so when defers have been generated then a local call is added to
// the return statement(s) for the block.
func LocalCallImpl(c *Context, i interface{}) error {

	// Make a new symbol table for the fucntion to run with,
	// and a new execution context. Store the argument list in
	// the child table.
	sf := symbols.NewChildSymbolTable("defer", c.symbols)
	cx := NewContext(sf, c.bc)
	cx.Tracing = c.Tracing

	cx.SetTokenizer(c.GetTokenizer())
	cx.result = nil

	// Make the caller's stack our stack
	cx.stack = c.stack
	cx.sp = c.sp

	// Run the function. If it doesn't get an error, then
	// extract the top stack item as the result
	err := cx.RunFromAddress(util.GetInt(i))

	// Because we share a stack with our caller, make sure the
	// caller's stack pointer is updated to match our value.
	c.sp = cx.sp
	return err

}

// CallImpl instruction processor calls a function (which can have
// parameters and a return value). The function value must be on the
// stack, preceded by the function arguments. The operand indicates the
// number of arguments that are on the stack. The function value must be
// etiher a pointer to a built-in function, or a pointer to a bytecode
// function implementation.
func CallImpl(c *Context, i interface{}) error {

	var err error
	var funcPointer interface{}

	// Argument count is in operand. It can be offset by a
	// value held in the context cause during argument processing.
	// Normally, this value is zero.
	argc := i.(int) + c.argCountDelta
	c.argCountDelta = 0

	// Arguments are in reverse order on stack.
	args := make([]interface{}, argc)
	for n := 0; n < argc; n = n + 1 {
		v, err := c.Pop()
		if err != nil {
			return err
		}
		args[(argc-n)-1] = v
	}

	// Function value is last item on stack
	funcPointer, err = c.Pop()
	if err != nil {
		return err
	}
	var result interface{}

	// Depends on the type here as to what we call...
	switch af := funcPointer.(type) {
	case *ByteCode:

		// Find the top of this scope level (typically)
		parentTable := c.symbols
		if !c.fullSymbolScope {
			for !parentTable.ScopeBoundary && parentTable.Parent != nil {
				parentTable = parentTable.Parent
			}
		}

		funcSymbols := symbols.NewChildSymbolTable("Function", parentTable)
		funcSymbols.ScopeBoundary = true

		// Make a new symbol table for the fucntion to run with,
		// and a new execution context. Note that this table has no
		// visibility into the current scope of symbol values.
		sf := symbols.NewChildSymbolTable("Function", parentTable)
		cx := NewContext(sf, af)
		cx.fullSymbolScope = c.fullSymbolScope
		cx.Tracing = c.Tracing
		cx.SetTokenizer(c.GetTokenizer())
		cx.result = nil

		// Make the caller's stack our stack
		cx.stack = c.stack
		cx.sp = c.sp

		_ = sf.SetAlways("_args", args)
		if c.this != nil {
			_ = sf.SetAlways("_this", c.this)
			c.this = nil
		}

		// Run the function. If it doesn't get an error, then
		// extract the top stack item as the result
		err = cx.Run()
		if err == nil {
			result = cx.result
		}

		// Because we share a stack with our caller, make sure the
		// caller's stack pointer is updated to match our value.
		c.sp = cx.sp

	case func(*symbols.SymbolTable, []interface{}) (interface{}, error):

		// First, can we check the argument count on behalf of the caller?
		df := functions.FindFunction(af)
		fname := runtime.FuncForPC(reflect.ValueOf(af).Pointer()).Name()
		fname = strings.Replace(fname, "github.com/tucats/gopackages/", "", 1)

		isUtilSymbols := (fname == "functions.FormatSymbols")

		if df != nil {
			if len(args) < df.Min || len(args) > df.Max {
				name := functions.FindName(af)
				return functions.NewError(name, ArgumentCountError)
			}
		}

		// Note special exclusion for the case of the util.Symbols function which must be
		// able to see the entire tree...
		parentTable := c.symbols
		if !isUtilSymbols && !c.fullSymbolScope {
			for !parentTable.ScopeBoundary && parentTable.Parent != nil {
				parentTable = parentTable.Parent
			}
		}

		funcSymbols := symbols.NewChildSymbolTable(fname, parentTable)
		funcSymbols.ScopeBoundary = true

		if c.this != nil {
			_ = funcSymbols.SetAlways("_this", c.this)
			c.this = nil
		}
		result, err = af(funcSymbols, args)

		if r, ok := result.(functions.MultiValueReturn); ok {
			_ = c.Push(StackMarker{Desc: "multivalue result"})
			for i := len(r.Value) - 1; i >= 0; i = i - 1 {
				_ = c.Push(r.Value[i])
			}
			return nil
		}
		// If there was an error but this function allows it, then
		// just push the result values
		if df != nil && df.ErrReturn {
			_ = c.Push(StackMarker{Desc: "builtin result"})
			_ = c.Push(err)
			_ = c.Push(result)
			return nil
		}

		// Functions implemented natively cannot wrap them up as runtime
		// errors, so let's help them out.
		if err != nil {
			name := functions.FindName(af)
			if name != "" {
				name = " " + name
			}
			err = c.NewError("in function" + name + ", " + err.Error())
		}

	default:
		return c.NewError(InvalidFunctionCallError, fmt.Sprintf("%#v", af))
	}

	if err != nil {
		return err
	}
	if result != nil {
		_ = c.Push(result)
	}
	return nil
}

// ReturnImpl implements the return opcode which returns from a called function
// or local subroutine
func ReturnImpl(c *Context, i interface{}) error {
	var err error
	// Do we have a return value?
	if b, ok := i.(bool); ok && b {
		c.result, err = c.Pop()
	}
	// Stop running this context
	c.running = false
	return err
}

// ArgCheckImpl instruction processor verifies that there are enough items
// on the stack to satisfy the function's argument list. The operand is the
// number of values that must be available. Alternaitvely, the operand can be
// an array of objects, which are the minimum count, maximum count, and
// function name
func ArgCheckImpl(c *Context, i interface{}) error {

	min := 0
	max := 0
	name := "function call"

	switch v := i.(type) {
	case []interface{}:
		if len(v) < 2 || len(v) > 3 {
			return c.NewError(InvalidArgCheckError)
		}
		min = util.GetInt(v[0])
		max = util.GetInt(v[1])
		if len(v) == 3 {
			name = util.GetString(v[2])
		}
	case int:
		if v >= 0 {
			min = v
			max = v
		} else {
			min = 0
			max = -v
		}

	case []int:
		if len(v) != 2 {
			return c.NewError(InvalidArgCheckError)
		}
		min = v[0]
		max = v[1]

	default:
		return c.NewError(InvalidArgCheckError)
	}

	v, found := c.Get("_args")
	if !found {
		return c.NewError(InvalidArgCheckError)
	}

	// Was there a "This" done just before this? If so, set
	// the stack value accordingly.
	if thisName, ok := c.this.(string); ok && thisName != "" {
		this, err := c.Pop()
		if err != nil {
			return err
		}
		_ = c.SetAlways(thisName, this)
		c.this = nil
	}

	// Do the actual compare. Note that if we ended up with a negative
	// max, that means variable argument list size, and we just assume
	// what we found in the max...
	va := v.([]interface{})
	if max < 0 {
		max = len(va)
	}
	if len(va) < min || len(va) > max {
		return functions.NewError(name, ArgumentCountError)
	}
	return nil
}

// TryImpl instruction processor
func TryImpl(c *Context, i interface{}) error {
	addr := util.GetInt(i)
	c.try = append(c.try, addr)
	return nil
}

// TryPopImpl instruction processor
func TryPopImpl(c *Context, i interface{}) error {
	if len(c.try) == 0 {
		return c.NewError(TryCatchMismatchError)
	}
	if len(c.try) == 1 {
		c.try = make([]int, 0)
	} else {
		c.try = c.try[:len(c.try)-1]
	}

	_ = c.symbols.DeleteAlways("_error")
	return nil
}
