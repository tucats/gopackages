package bytecode

import (
	"fmt"

	"github.com/tucats/gopackages/functions"
	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/util"
)

/******************************************\
*                                         *
*        F L O W   C O N T R O L          *
*                                         *
\******************************************/

// StopOpcode bytecode implementation
func StopOpcode(c *Context, i interface{}) error {
	c.running = false
	return nil
}

// PanicOpcode bytecode implementation. The boolean
// flag has to indicate if this is a fatal error
func PanicOpcode(c *Context, i interface{}) error {
	c.running = !util.GetBool(i)
	strValue, err := c.Pop()
	if err != nil {
		return err
	}
	msg := util.GetString(strValue)
	return c.NewError(msg)
}

// AtLineOpcode implementation. This identifies the
// start of a new statement, and tags the line number
// from the source where this was found. This is used
// in error messaging, primarily.
func AtLineOpcode(c *Context, i interface{}) error {

	c.line = util.GetInt(i)
	// If we are tracing, put that out now.
	if c.tokenizer != nil {
		fmt.Printf("%d:  %s\n", c.line, c.tokenizer.GetLine(c.line))
	}
	return nil
}

// BranchFalseOpcode bytecode implementation
func BranchFalseOpcode(c *Context, i interface{}) error {

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

// BranchOpcode bytecode implementation
func BranchOpcode(c *Context, i interface{}) error {

	// Get destination
	address := util.GetInt(i)
	if address < 0 || address > c.bc.emitPos {
		return c.NewError(InvalidBytecodeAddress)
	}

	c.pc = address
	return nil
}

// BranchTrueOpcode bytecode implementation
func BranchTrueOpcode(c *Context, i interface{}) error {

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

func LocalCallOpcode(c *Context, i interface{}) error {

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

// CallOpcode bytecode implementation.
func CallOpcode(c *Context, i interface{}) error {

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

		// Make a new symbol table for the fucntion to run with,
		// and a new execution context. Store the argument list in
		// the child table.
		sf := symbols.NewChildSymbolTable("Function", c.symbols)
		cx := NewContext(sf, af)
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
		if df != nil {
			if len(args) < df.Min || len(args) > df.Max {
				name := functions.FindName(af)
				return functions.NewError(name, ArgumentCountError)
			}
		}
		if c.this != nil {
			_ = c.symbols.SetAlways("_this", c.this)
			c.this = nil
		}

		result, err = af(c.symbols, args)

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

// ReturnOpcode implements the return opcode which returns
// from a called function.
func ReturnOpcode(c *Context, i interface{}) error {
	var err error
	// Do we have a return value?
	if b, ok := i.(bool); ok && b {
		c.result, err = c.Pop()
	}
	// Stop running this context
	c.running = false
	return err
}

// ArgCheckOpcode implementation
func ArgCheckOpcode(c *Context, i interface{}) error {

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

// TryOpcode implementation
func TryOpcode(c *Context, i interface{}) error {
	addr := util.GetInt(i)
	c.try = append(c.try, addr)
	return nil
}

// TryPopOpcode implementation
func TryPopOpcode(c *Context, i interface{}) error {
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
