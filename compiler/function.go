package compiler

import (
	"fmt"

	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/tokenizer"
	"github.com/tucats/gopackages/util"
)

// Function compiles a function definition. The literal flag indicates if
// this is a function literal, which is pushed on the stack, or a non-literal
// which is added to the symbol table dictionary.
func (c *Compiler) Function(literal bool) error {

	// List of type coercions that will be needed for any RETURN statement.
	coercions := []*bytecode.ByteCode{}

	// Descriptor of each parameter in the parameter list.
	type parameter struct {
		name string
		kind int
	}
	parameters := []parameter{}
	this := ""
	fname := ""

	// If it's not a literal, there will be a function name, which must be a valid
	// symbol name. It might also be an object-oriented (a->b()) call.
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

	// Process the function parameter specification
	varargs := false
	if c.t.IsNext("(") {
		for !c.t.IsNext(")") {
			if c.t.AtEnd() {
				break
			}
			name := c.t.Next()
			p := parameter{kind: bytecode.UndefinedType}
			if tokenizer.IsSymbol(name) {
				p.name = name
			} else {
				return c.NewError(InvalidFunctionArgument)
			}
			if c.t.Peek(1) == "..." {
				c.t.Advance(1)
				varargs = true
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
			} else if util.InList(c.t.Peek(1), "interface{}", "int", "string", "bool", "double", "float", "array", "struct") {
				switch c.t.Next() {
				case "int":
					p.kind = bytecode.IntType
				case "string":
					p.kind = bytecode.StringType
				case "bool":
					p.kind = bytecode.BoolType
				case "float", "double":
					p.kind = bytecode.FloatType
				case "struct":
					p.kind = bytecode.StructType
				case "array":
					p.kind = bytecode.ArrayType
				}
			}
			if varargs {
				p.kind = bytecode.VarArgs
			}
			parameters = append(parameters, p)
			_ = c.t.IsNext(",")
		}
	}

	b := bytecode.New(fname)

	// If we know our source file, mark it in the bytecode now.
	if c.SourceFile != "" {
		b.Emit(bytecode.FromFile, c.SourceFile)
	}

	// Generate the argument check
	p := []interface{}{
		len(parameters),
		len(parameters),
		fname,
	}
	if varargs {
		p[0] = len(parameters)
		p[1] = -1
	}

	b.Emit(bytecode.AtLine, c.t.Line[c.t.TokenP])
	b.Emit(bytecode.ArgCheck, p)

	// If there was a "this" variable defined, process it now.
	if this != "" {
		b.Emit(bytecode.This, this)
	}

	// Generate the parameter assignments. These are extracted from the automatic
	// array named _args which is generated as part of the bytecode function call.
	for n, p := range parameters {

		// is this the end of the fixed list? If so, emit the instruction that scoops
		// up the remaining arguments and stores them as an array value.
		//
		// Otherwise, generate code to extract the argument value by index number.
		if p.kind == bytecode.VarArgs {
			b.Emit(bytecode.GetVarArgs, n)
		} else {
			b.Emit(bytecode.Load, "__args")
			b.Emit(bytecode.Push, n)
			b.Emit(bytecode.LoadIndex)
		}

		// If this argumnet is not interface{} or a variable argument item,
		// generaet code to validate/coerce the value to a given type.
		if p.kind != bytecode.UndefinedType && p.kind != bytecode.VarArgs {
			b.Emit(bytecode.RequiredType, p.kind)
		}
		// Generate code to store the value on top of the stack into the local
		// symbol for the parameter name.
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
				case "float", "double":
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
				case "interface{}":
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
	// current token stream in progress, and the current bytecode. But otherwise we
	// use a new compiler context, so any nested operations do not affect the definition
	// of the function body we're compiling.
	cx := New()
	cx.t = c.t
	cx.b = b
	cx.coerce = coercions
	err := cx.Statement()
	if err != nil {
		return err
	}
	// Add trailing return to ensure we close out the scope correctly
	cx.b.Emit(bytecode.Return)

	// If it was a literal, push the body of the function (really, a bytecode expression
	// of the function code) on the stack. Otherwise, let's store it in the symbol table
	// or package dictionary as appropraite.
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
