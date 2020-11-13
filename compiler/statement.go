package compiler

import (
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/tokenizer"
)

// Predefined names used by statement processing.
const (
	DirectiveStructureName = "_directives"
)

// Statement compiles a single statement
func (c *Compiler) Statement() error {

	// We just eat statement separators and we also terminate
	// processing when we hit the end of the token stream.
	if c.t.IsNext(";") {
		return nil
	}
	if c.t.IsNext(tokenizer.EndOfTokens) {
		return nil
	}

	// Is it a directive token? These really just store data in the compiler
	// symbol table that is used to extend features. These symbols end up in
	// the runtime context of the running code

	if c.t.IsNext("@") {
		return c.Directive()
	}
	c.statementCount = c.statementCount + 1

	// Is it a function definition? These aren't compiled inline,
	// so we call a special compile unit that will compile the
	// function and store it in the bytecode symbol table.
	if c.t.IsNext("func") {
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
	// which compiler unit to call. For each term, call the appropriate
	// handler (which assumes the leading verb has already been consumed)
	switch c.t.Next() {
	case "{":
		return c.Block()
	case "array":
		return c.Array()
	case "assert":
		return c.Assert()
	case "break":
		return c.Break()
	case "call":
		return c.Call()
	case "const":
		return c.Constant()
	case "continue":
		return c.Continue()
	case "exit":
		return c.Exit()
	case "for":
		return c.For()
	case "if":
		return c.If()
	case "import":
		return c.Import()
	case "package":
		return c.Package()
	case "print":
		return c.Print()
	case "return":
		return c.Return()
	case "switch":
		return c.Switch()
	case "try":
		return c.Try()
	case "type":
		return c.Type()
	}

	// Unknown statement, return an error
	return c.NewTokenError("unrecognized or unexpected token")
}
