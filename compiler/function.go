package compiler

import (
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/tokenizer"
)

// Function compiles a function definition
func (c *Compiler) Function() error {

	parameters := []string{}
	this := ""

	fname := c.t.Next()
	if !tokenizer.IsSymbol(fname) {
		return c.NewTokenError("invalid function name")
	}

	// Was it really the function name, or the "this" variable name?
	if c.t.Peek(1) == "->" {
		c.t.Advance(1)
		this = fname
		fname = c.t.Next()
		if !tokenizer.IsSymbol(fname) {
			return c.NewTokenError("invalid function name")
		}
	}

	// Process parameter names
	varargs := false
	if c.t.IsNext("(") {
		for !c.t.IsNext(")") {
			if c.t.AtEnd() {
				break
			}
			if c.t.Peek(1) == "." && c.t.Peek(2) == "." && c.t.Peek(3) == "." {
				c.t.Advance(3)
				varargs = true
				continue
			}
			name := c.t.Next()
			if tokenizer.IsSymbol(name) {
				parameters = append(parameters, name)
			} else {
				return c.NewTokenError("invalid parameter")
			}
			if c.t.IsNext(",") {
				// No action
			}
		}
	}

	b := bytecode.New(fname)

	// Generate the argument check
	if varargs {
		p := make([]int, 2)
		p[0] = len(parameters)
		p[1] = 999999
		b.Emit2(bytecode.ArgCheck, p)

	} else {
		b.Emit2(bytecode.ArgCheck, len(parameters))
	}

	// If there was a "this" variable defined, process it now.
	if this != "" {
		b.Emit2(bytecode.This, this)
	}

	// Generate the parameter assignments. These are extracted
	// from the automatic array named _args which is generated
	// as part of the function call during bytecode exectuion.
	// Note that the array is 1-based.
	for n, name := range parameters {
		b.Emit2(bytecode.Load, "_args")
		b.Emit2(bytecode.Push, n+1)
		b.Emit1(bytecode.LoadIndex)
		b.Emit2(bytecode.SymbolCreate, name)
		b.Emit2(bytecode.Store, name)
	}

	// Look for return type definition. If found, compile the appropriate
	// coercion code which will be stored in the comipler block for use
	// by a return statement

	coercion := bytecode.New(fname + " return")

	if c.t.Peek(1) == "[" && c.t.Peek(2) == "]" {
		coercion.Emit2(bytecode.Coerce, bytecode.ArrayType)
		c.t.Advance(2)
	} else {
		if c.t.Peek(1) == "{" && c.t.Peek(2) == "}" {
			coercion.Emit2(bytecode.Coerce, bytecode.StructType)
			c.t.Advance(2)
		} else {
			switch c.t.Peek(1) {
			case "int":
				coercion.Emit2(bytecode.Coerce, bytecode.IntType)
				c.t.Advance(1)
			case "float":
				coercion.Emit2(bytecode.Coerce, bytecode.FloatType)
				c.t.Advance(1)
			case "string":
				coercion.Emit2(bytecode.Coerce, bytecode.StringType)
				c.t.Advance(1)
			case "bool":
				coercion.Emit2(bytecode.Coerce, bytecode.BoolType)
				c.t.Advance(1)

			case "any":
				coercion.Emit2(bytecode.Coerce, bytecode.UndefinedType)
				c.t.Advance(1)

			case "void":
				// Do nothing, there is no result.
				c.t.Advance(1)

			default:
				return c.NewError("missing function return type")
			}
		}
	}
	// Now compile a statement or block into the function body. We'll use the
	// current token stream in progress, and the current bytecode.
	cx := New()
	cx.t = c.t
	cx.b = b
	cx.coerce = coercion
	err := cx.Statement()
	if err != nil {
		return err
	}

	// Store anchor to the function, either in the current
	// table or package.

	if c.PackageName == "" {
		c.s.SetAlways(fname, b)
	} else {
		c.AddPackageFunction(c.PackageName, fname, b)
	}

	return nil
}
