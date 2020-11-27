package bytecode

import (
	"fmt"
	"strconv"

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
		return c.NewError("invalid destination address: " + strconv.Itoa(address))
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
		return c.NewError("invalid destination address: " + strconv.Itoa(address))
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
		return c.NewError("invalid destination address: " + strconv.Itoa(address))
	}

	if util.GetBool(v) {
		c.pc = address
	}
	return nil
}

// CallOpcode bytecode implementation.
func CallOpcode(c *Context, i interface{}) error {

	var err error
	var funcPointer interface{}

	// Argument count is in operand
	argc := i.(int)

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

		sf.SetAlways("_args", args)
		if c.this != nil {
			sf.SetAlways("_this", c.this)
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
				if name > "" {
					name = ": " + name
				}
				if len(args) < df.Min {
					return c.NewError("insufficient arguments" + name)
				}
				return c.NewError("too many arguments" + name)
			}
		}
		if c.this != nil {
			c.symbols.SetAlways("_this", c.this)
			c.this = nil
		}

		result, err = af(c.symbols, args)

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
		return c.NewStringError("invalid target of function call", fmt.Sprintf("%#v", af))
	}

	if err != nil {
		return err
	}
	if result != nil {
		c.Push(result)
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

	switch v := i.(type) {
	case []interface{}:
		if len(v) != 2 {
			return c.NewError("invalid ArgCheck array size")
		}
		min = v[0].(int)
		max = v[1].(int)

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
			return c.NewError("invalid ArgCheck array size")
		}
		min = v[0]
		max = v[1]

	default:
		return c.NewError("invalid ArgCheck operand")
	}

	v, found := c.Get("_args")
	if !found {
		return c.NewError("ArgCheck cannot read _args")
	}

	// Was there a "This" done just before this? If so, set
	// the stack value accordingly.
	if thisName, ok := c.this.(string); ok && thisName != "" {
		this, err := c.Pop()
		if err != nil {
			return err
		}
		c.SetAlways(thisName, this)
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
		return c.NewError("incorrect number of arguments passed")
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
		return c.NewError("try/catch mismatch")
	}
	if len(c.try) == 1 {
		c.try = make([]int, 0)
	} else {
		c.try = c.try[:len(c.try)-1]
	}

	c.symbols.DeleteAlways("_error")
	return nil
}
