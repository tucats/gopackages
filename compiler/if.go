package compiler

import (
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/expressions"
)

// If compiles conditional statments. The verb is already
// removed from the token stream.
func (c *Compiler) If() error {

	// Compile the conditional expression
	bc, err := expressions.Compile(c.t)
	if err != nil {
		return err
	}
	c.b.Append(bc)
	c.b.Emit(bytecode.Push, "bool")
	c.b.Emit(bytecode.Call, 1)

	b1 := c.b.Mark()
	c.b.Emit(bytecode.BranchFalse, 0)

	// Compile the statement to be executed if true

	err = c.Statement()
	if err != nil {
		return err
	}

	// If there's an else clause, branch around it.
	if c.t.IsNext("else") {
		b2 := c.b.Mark()
		c.b.Emit(bytecode.Branch, 0)
		c.b.SetAddressHere(b1)

		err = c.Statement()
		if err != nil {
			return err
		}
		c.b.SetAddressHere(b2)

	} else {
		c.b.SetAddressHere(b1)
	}
	return nil
}
