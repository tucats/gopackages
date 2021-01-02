package compiler

import (
	"fmt"

	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/tokenizer"
	"github.com/tucats/gopackages/util"
)

// Function compiles a function definition. If the literal flag is
// set, this generates a function and puts the pointer to it on the
// stack. IF not set, the function is added to the global or local
// function dictionary.
func (c *Compiler) Function(literal bool) error {

	coercions := []*bytecode.ByteCode{}

	type parameter struct {
		name string
		kind int
	}
	parameters := []parameter{}
	this := ""
	fname := ""

	if !literal {
		fname = c.t.Next()
		if !tokenizer.IsSymbol(fname) {
			return c.NewError(InvalidFunctionName, fname)
		}
		fname = c.Normalize(fname)

		// Was it really the function name, or the "this" variable name?
		if c.t.Peek(1) == "->" {
			c.t.Advance(1)
			this = fname
			fname = c.t.Next()
			if !tokenizer.IsSymbol(fname) {
				return c.NewError(InvalidFunctionName, fname)
			}
			fname = c.Normalize(fname)
		}
	}

	// Process parameter names
	varargs := false
	if c.t.IsNext("(") {
		for !c.t.IsNext(")") {
			if c.t.AtEnd() {
				break
			}
			if c.t.Peek(1) == "..." {
				c.t.Advance(1)
				varargs = true
				continue
			}
			name := c.t.Next()
			p := parameter{kind: bytecode.UndefinedType}
			if tokenizer.IsSymbol(name) {
				p.name = name
			} else {
				return c.NewError(InvalidFunctionArgument)
			}

			// Is there a type name that follows it? We have to check for "[]" and "{}"
			// as two differnt tokens. Also note that you can use the word array or struct
			// instead if you wish.
			if c.t.Peek(1) == "[" && c.t.Peek(2) == "]" {
				p.kind = bytecode.ArrayType
				c.t.Advance(2)
			} else if c.t.Peek(1) == "{" && c.t.Peek(2) == "}" {
				p.kind = bytecode.ArrayType
				c.t.Advance(2)
			} else if util.InList(c.t.Peek(1), "any", "int", "string", "bool", "float", "array", "struct") {
				switch c.t.Next() {
				case "int":
					p.kind = bytecode.IntType
				case "string":
					p.kind = bytecode.StringType
				case "bool":
					p.kind = bytecode.BoolType
				case "float":
					p.kind = bytecode.FloatType
				case "struct":
					p.kind = bytecode.StructType
				case "array":
					p.kind = bytecode.ArrayType
				}
			}

			parameters = append(parameters, p)
			_ = c.t.IsNext(",")
		}
	}

	b := bytecode.New(fname)

	// Generate the argument check
	p := []interface{}{
		len(parameters),
		len(parameters),
		fname,
	}
	if varargs {
		p[1] = -1
	}
	b.Emit(bytecode.ArgCheck, p)

	// If there was a "this" variable defined, process it now.
	if this != "" {
		b.Emit(bytecode.This, this)
	}

	// Generate the parameter assignments. These are extracted
	// from the automatic array named _args which is generated
	// as part of the function call during bytecode execution.
	for n, p := range parameters {
		b.Emit(bytecode.Load, "_args")
		b.Emit(bytecode.Push, n)
		b.Emit(bytecode.LoadIndex)
		if p.kind != bytecode.UndefinedType {
			b.Emit(bytecode.Coerce, p.kind)
		}
		b.Emit(bytecode.SymbolCreate, p.name)
		b.Emit(bytecode.Store, p.name)
	}

	// Is there a list of return items (expressed as a parenthesis)?
	hasReturnList := c.t.IsNext("(")
	returnValueCount := 0
	wasVoid := false
	// Loop over the (possibly singular) return type specification
	for {
		coercion := bytecode.New(fmt.Sprintf("%s return item %d", fname, returnValueCount))
		if c.t.Peek(1) == "[" && c.t.Peek(2) == "]" {
			coercion.Emit(bytecode.Coerce, bytecode.ArrayType)
			c.t.Advance(2)
		} else {
			if c.t.Peek(1) == "{" && c.t.Peek(2) == "}" {
				coercion.Emit(bytecode.Coerce, bytecode.StructType)
				c.t.Advance(2)
			} else {
				switch c.t.Peek(1) {
				// Start of block means no more types here.
				case "{":
					break
				case "error":
					c.t.Advance(1)
					coercion.Emit(bytecode.Coerce, bytecode.ErrorType)

				case "int":
					coercion.Emit(bytecode.Coerce, bytecode.IntType)
					c.t.Advance(1)
				case "float":
					coercion.Emit(bytecode.Coerce, bytecode.FloatType)
					c.t.Advance(1)
				case "string":
					coercion.Emit(bytecode.Coerce, bytecode.StringType)
					c.t.Advance(1)
				case "bool":
					coercion.Emit(bytecode.Coerce, bytecode.BoolType)
					c.t.Advance(1)
				case "struct":
					coercion.Emit(bytecode.Coerce, bytecode.StructType)
					c.t.Advance(1)
				case "array":
					coercion.Emit(bytecode.Coerce, bytecode.ArrayType)
					c.t.Advance(1)
				case "any":
					coercion.Emit(bytecode.Coerce, bytecode.UndefinedType)
					c.t.Advance(1)

				case "void":
					// Do nothing, there is no result.
					wasVoid = true
					c.t.Advance(1)

				default:
					return c.NewError(MissingFunctionTypeError)
				}
			}
		}
		if !wasVoid {
			coercions = append(coercions, coercion)
		}

		if c.t.Peek(1) != "," {
			break
		}
		// If we got here, but never had a () around this list, it's an error
		if !hasReturnList {
			return c.NewError(InvalidReturnTypeList)
		}
		c.t.Advance(1)
	}

	// If the return types were expressed as a list, there must be a trailing paren.
	if hasReturnList && !c.t.IsNext(")") {
		return c.NewError(MissingParenthesisError)
	}

	// Now compile a statement or block into the function body. We'll use the
	// current token stream in progress, and the current bytecode.
	cx := New()
	cx.t = c.t
	cx.b = b
	cx.coerce = coercions
	err := cx.Statement()
	if err != nil {
		return err
	}

	if literal {
		c.b.Emit(bytecode.Push, b)
	} else {
		// Store address of the function, either in the current
		// compiler's symbol table or active package.
		if c.PackageName == "" {
			_ = c.s.SetAlways(fname, b)
		} else {
			_ = c.AddPackageFunction(c.PackageName, fname, b)
		}
	}
	return nil
}
