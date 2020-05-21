package compiler

import (
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/tokenizer"
)

// Statement parses a single statement
func (c *Compiler) Statement() error {

	// We just eat statement separators and we also terminate
	// processing when we hit the end of the token stream.
	if c.t.IsNext(";") {
		return nil
	}
	if c.t.IsNext(tokenizer.EndOfTokens) {
		return nil
	}

	// Is it a function definition? These aren't compiled inline,
	// so we call a special compile unit that will compile the
	// function and store it in the bytecode symbol table.
	if c.t.IsNext("function") {
		return c.Function()
	}

	// At this point, we know we're trying to compile a statement,
	// so store the current line number in the stream to help us
	// form runtime error messages as needed.
	c.b.Emit2(bytecode.AtLine, c.t.Line[c.t.TokenP])

	// If the next item(s) constitute a value LValue, then this is
	// an assignment statement.
	if c.IsLValue() {
		return c.Assignment()
	}

	// Remaining statement types all have a starting term that defines
	// which compiler unit to call.
	switch c.t.Next() {

	case "package":
		return c.Package()

	case "import":
		return c.Import()

	case "const":
		return c.Constant()

	case "{":
		return c.Block()

	case "if":
		return c.If()

	case "for":
		return c.For()

	case "break":
		return c.Break()

	case "continue":
		return c.Continue()

	case "print":
		return c.Print()

	case "call":
		return c.Call()

	case "return":
		return c.Return()

	case "array":
		return c.Array()

	}

	// Unknown statement, return an error
	return c.NewTokenError("unrecognized or unexpected token")
}
