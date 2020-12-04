package compiler

import (
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/tokenizer"
)

// Function compiles a function definition. If the literal flag is
// set, this generates a function and puts the pointer to it on the
// stack. IF not set, the function is added to the global or local
// function dictionary.
func (c *Compiler) Function(literal bool) error {

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

		// Was it really the function name, or the "this" variable name?
		if c.t.Peek(1) == "->" {
			c.t.Advance(1)
			this = fname
			fname = c.t.Next()
			if !tokenizer.IsSymbol(fname) {
				return c.NewError(InvalidFunctionName, fname)
			}
		}
	}
	name = c.Normalize(name)
	fname = c.Normalize(fname)

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
			p := parameter{kind: bytecode.UndefinedType}
			if tokenizer.IsSymbol(name) {
				p.name = name
			} else {
				return c.NewError(InvalidFunctionArgument)
			}
			name = c.Normalize(name)

			// Is there a type name that follows it? We have to check for "[]" and "{}"
			// as two differnt tokens. Also note that you can use the word array or struct
			// instead if you wish.
			if c.t.Peek(1) == "[" && c.t.Peek(2) == "]" {
				p.kind = bytecode.ArrayType
				c.t.Advance(2)
			} else if c.t.Peek(1) == "{" && c.t.Peek(2) == "}" {
				p.kind = bytecode.ArrayType
				c.t.Advance(2)
			} else if inList(c.t.Peek(1), []string{"any", "int", "string", "bool", "float", "array", "struct"}) {
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
			if c.t.IsNext(",") {
				// No action
			}
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

	// Look for return type definition. If found, compile the appropriate
	// coercion code which will be stored in the compiler block for use
	// by a return statement
	coercion := bytecode.New(fname + " return")

	if c.t.Peek(1) == "[" && c.t.Peek(2) == "]" {
		coercion.Emit(bytecode.Coerce, bytecode.ArrayType)
		c.t.Advance(2)
	} else {
		if c.t.Peek(1) == "{" && c.t.Peek(2) == "}" {
			coercion.Emit(bytecode.Coerce, bytecode.StructType)
			c.t.Advance(2)
		} else {
			switch c.t.Peek(1) {
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
				c.t.Advance(1)

			default:
				return c.NewError(MissingFunctionTypeError)
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

	if literal {
		c.b.Emit(bytecode.Push, b)
	} else {
		// Store address of the function, either in the current
		// compiler's symbol table or active package.
		if c.PackageName == "" {
			c.s.SetAlways(fname, b)
		} else {
			c.AddPackageFunction(c.PackageName, fname, b)
		}
	}
	return nil
}

func inList(search string, values []string) bool {
	for _, item := range values {
		if search == item {
			return true
		}
	}
	return false
}
