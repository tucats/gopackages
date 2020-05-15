package compiler

import (
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/expressions"
)

// Call handles the call statement. This is really the same as
// invoking a function in an expression, except there is no
// result value.
func (c *Compiler) Call() error {

	// Let's peek ahead to see if this is a legit function call

	if !expressions.Symbol(c.t.Peek(1)) || (c.t.Peek(2) != "(" && c.t.Peek(2) != ".") {
		return c.NewError("invalid function call")
	}

	// Parse the function as an expression, which we then ignore the
	// result of.
	bc, err := expressions.Compile(c.t)
	if err != nil {
		return err
	}
	c.b.Append(bc)
	c.b.Emit0(bytecode.Drop)
	return nil
}
