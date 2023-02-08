package bytecode

import (
	"reflect"
	"runtime"
	"strings"

	"github.com/tucats/gopackages/defs"
	"github.com/tucats/gopackages/errors"
	"github.com/tucats/gopackages/expressions/builtins"
	"github.com/tucats/gopackages/expressions/data"
	"github.com/tucats/gopackages/expressions/symbols"
)

/******************************************\
*                                         *
*        F L O W   C O N T R O L          *
*                                         *
\******************************************/

// stopByteCode instruction processor causes the current execution context to
// stop executing immediately.
func stopByteCode(c *Context, i interface{}) error {
	c.running = false

	return errors.ErrStop
}

// branchFalseByteCode instruction processor branches to the instruction named in
// the operand if the top-of-stack item is a boolean FALSE value. Otherwise,
// execution continues with the next instruction.
func branchFalseByteCode(c *Context, i interface{}) error {
	// Get test value
	v, err := c.Pop()
	if err != nil {
		return err
	}

	if c.typeStrictness == 0 {
		if _, ok := v.(bool); !ok {
			return c.error(errors.ErrConditionalBool).Context(data.TypeOf(v).String())
		}
	}

	// Get destination

	if address := data.Int(i); address < 0 || address > c.bc.nextAddress {
		return c.error(errors.ErrInvalidBytecodeAddress).Context(address)
	} else {
		if !data.Bool(v) {
			c.programCounter = address
		}
	}

	return nil
}

// branchByteCode instruction processor branches to the instruction named in
// the operand.
func branchByteCode(c *Context, i interface{}) error {
	// Get destination
	if address := data.Int(i); address < 0 || address > c.bc.nextAddress {
		return c.error(errors.ErrInvalidBytecodeAddress).Context(address)
	} else {
		c.programCounter = address
	}

	return nil
}

// branchTrueByteCode instruction processor branches to the instruction named in
// the operand if the top-of-stack item is a boolean TRUE value. Otherwise,
// execution continues with the next instruction.
func branchTrueByteCode(c *Context, i interface{}) error {
	// Get test value
	v, err := c.Pop()
	if err != nil {
		return err
	}

	if c.typeStrictness == 0 {
		if _, ok := v.(bool); !ok {
			return c.error(errors.ErrConditionalBool).Context(data.TypeOf(v).String())
		}
	}

	// Get destination
	if address := data.Int(i); address < 0 || address > c.bc.nextAddress {
		return c.error(errors.ErrInvalidBytecodeAddress).Context(address)
	} else {
		if data.Bool(v) {
			c.programCounter = address
		}
	}

	return nil
}

// callByteCode instruction processor calls a function (which can have
// parameters and a return value). The function value must be on the
// stack, preceded by the function arguments. The operand indicates the
// number of arguments that are on the stack. The function value must be
// either a pointer to a built-in function, or a pointer to a bytecode
// function implementation.
func callByteCode(c *Context, i interface{}) error {
	var err error

	var functionPointer interface{}

	var result interface{}

	// Argument count is in operand. It can be offset by a
	// value held in the context cause during argument processing.
	// Normally, this value is zero.
	argc := data.Int(i) + c.argCountDelta
	c.argCountDelta = 0
	fullSymbolVisibility := c.fullSymbolScope

	// Determine if language extensions are supported. This is required
	// for variable length argument lists that are not variadic.
	extensions := false

	if v, found := c.symbols.Get(defs.ExtensionsVariable); found {
		extensions = data.Bool(v)
	}

	// Arguments are in reverse order on stack.
	args := make([]interface{}, argc)

	for n := 0; n < argc; n = n + 1 {
		v, err := c.Pop()
		if err != nil {
			return err
		}

		if isStackMarker(v) {
			return c.error(errors.ErrFunctionReturnedVoid)
		}

		args[(argc-n)-1] = v
	}

	// Function value is last item on stack
	functionPointer, err = c.Pop()
	if err != nil {
		return err
	}

	if functionPointer == nil {
		return c.error(errors.ErrInvalidFunctionCall).Context("<nil>")
	}

	if isStackMarker(functionPointer) {
		return c.error(errors.ErrFunctionReturnedVoid)
	}

	// If this is a function pointer (from a stored type function list)
	// unwrap the value of the function pointer.
	if dp, ok := functionPointer.(data.Function); ok {
		fargc := 0

		if dp.Declaration != nil {
			fargc = len(dp.Declaration.Parameters)
			fullSymbolVisibility = dp.Declaration.Scope
		}

		if fargc != argc {
			// If extensions are not enabled, we don't allow variable argument counts.
			if !extensions && dp.Declaration != nil && !dp.Declaration.Variadic {
				return c.error(errors.ErrArgumentCount)
			}

			if fargc > 0 && (dp.Declaration.ArgCount[0] != 0 || dp.Declaration.ArgCount[1] != 0) {
				if argc < dp.Declaration.ArgCount[0] || argc > dp.Declaration.ArgCount[1] {
					return c.error(errors.ErrArgumentCount)
				}
			}
		}

		if c.typeStrictness == defs.StrictTypeEnforcement && dp.Declaration != nil {
			for n, arg := range args {
				parms := dp.Declaration.Parameters

				if dp.Declaration.Variadic && n > len(parms) {
					lastType := dp.Declaration.Parameters[len(parms)-1].Type

					if lastType.IsInterface() || lastType.IsType(data.ArrayType(data.InterfaceType)) || lastType.IsType(data.PointerType(data.InterfaceType)) {
						continue
					}

					if !data.TypeOf(arg).IsType(lastType) {
						return c.error(errors.ErrArgumentType).Context(data.TypeOf(arg).String())
					}
				}

				if n < len(parms) {
					if parms[n].Type.IsInterface() {
						continue
					}

					if parms[n].Type.IsType(data.ArrayType(data.InterfaceType)) || parms[n].Type.IsType(data.PointerType(data.InterfaceType)) {
						continue
					}

					if data.TypeOf(arg).IsInterface() {
						continue
					}

					if !data.TypeOf(arg).IsType(parms[n].Type) {
						return c.error(errors.ErrArgumentType).Context(data.TypeOf(arg).String())
					}
				}
			}
		}

		functionPointer = dp.Value
	}

	// Depends on the type here as to what we call...
	switch function := functionPointer.(type) {
	case *data.Type:
		// Calls to a type are really an attempt to cast the value.
		args = append(args, function)

		v, err := builtins.Cast(c.symbols, args)
		if err == nil {
			err = c.push(v)
		}

		return err

	case builtins.NativeFunction:
		// Native functions are methods on actual Go objects that we surface to Ego
		// code. Examples include the functions for waitgroup and mutex objects.
		functionName := builtins.GetName(function)
		funcSymbols := symbols.NewChildSymbolTable("builtin "+functionName, c.symbols)

		if v, ok := c.popThis(); ok {
			funcSymbols.SetAlways(defs.ThisVariable, v)
		}

		result, err = function(funcSymbols, args)

		if r, ok := result.(data.Values); ok {
			_ = c.push(NewStackMarker("results", 0))
			for i := len(r.Items) - 1; i >= 0; i = i - 1 {
				_ = c.push(r.Items[i])
			}

			return nil
		}

		// Functions implemented natively cannot wrap them up as runtime
		// errors, so let's help them out.
		if err != nil {
			err = c.error(err).In(builtins.FindName(function))
		}

	case func(*symbols.SymbolTable, []interface{}) (interface{}, error):
		// First, can we check the argument count on behalf of the caller?
		functionDefinition := builtins.FindFunction(function)
		functionName := runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name()
		functionName = strings.Replace(functionName, "github.com/tucats/gopackages/", "", 1)

		// See if it is a builtin function that needs visibility to the entire
		// symbol stack without binding the scope to the parent of the current
		// stack.
		if functionDefinition != nil {
			fullSymbolVisibility = fullSymbolVisibility || functionDefinition.FullScope

			if len(args) < functionDefinition.Min || len(args) > functionDefinition.Max {
				name := builtins.FindName(function)

				return c.error(errors.ErrArgumentCount).Context(name)
			}
		}

		// Note special exclusion for the case of the util.Symbols function which must be
		// able to see the entire tree...
		parentTable := c.symbols

		if !fullSymbolVisibility {
			for !parentTable.ScopeBoundary() && parentTable.Parent() != nil {
				parentTable = parentTable.Parent()
			}
		}

		functionSymbols := symbols.NewChildSymbolTable("builtin "+functionName, parentTable)
		functionSymbols.SetScopeBoundary(true)

		// Is this builtin one that requires a "this" variable? If so, get it from
		// the "this" stack.
		if v, ok := c.popThis(); ok {
			functionSymbols.SetAlways(defs.ThisVariable, v)
		}

		result, err = function(functionSymbols, args)

		if results, ok := result.(data.Values); ok {
			_ = c.push(NewStackMarker("results", 0))

			for i := len(results.Items) - 1; i >= 0; i = i - 1 {
				_ = c.push(results.Items[i])
			}

			return nil
		}

		// If there was an error but this function allows it, then
		// just push the result values
		if functionDefinition != nil && functionDefinition.ErrReturn {
			_ = c.push(NewStackMarker("results", 0))
			_ = c.push(err)
			_ = c.push(result)

			return nil
		}

		// Functions implemented natively cannot wrap them up as runtime
		// errors, so let's help them out.
		if err != nil {
			err = c.error(err)
		}

	case error:
		return c.error(errors.ErrUnusedErrorReturn)

	default:
		return c.error(errors.ErrInvalidFunctionCall).Context(function)
	}

	// IF no problems and there's a result value, push it on the
	// stack now.
	if err == nil && result != nil {
		err = c.push(result)
	}

	return err
}
