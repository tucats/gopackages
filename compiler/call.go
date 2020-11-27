package compiler

import (
	"github.com/tucats/gopackages/expressions"
	"github.com/tucats/gopackages/tokenizer"
)

// Call handles the call statement. This is really the same as
// invoking a function in an expression, except there is no
// result value.
func (c *Compiler) Call() error {

	// Let's peek ahead to see if this is a legit function call
	if !tokenizer.IsSymbol(c.t.Peek(1)) || (c.t.Peek(2) != "->" && c.t.Peek(2) != "(" && c.t.Peek(2) != ".") {
		return c.NewError(InvalidFunctionCall)
	}

	// Parse the function as an expression, which we then ignore the
	// result of.
	bc, err := expressions.Compile(c.t)
	if err != nil {
		return err
	}
	c.b.Append(bc)
	return nil
}
