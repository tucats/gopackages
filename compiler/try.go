package compiler

import "github.com/tucats/gopackages/bytecode"

// Try comiples the try statement
func (c *Compiler) Try() error {

	// Generate start of a try block.
	b1 := c.b.Mark()
	c.b.Emit2(bytecode.Try, 0)

	// Statement to try
	err := c.Statement()
	if err != nil {
		return err
	}
	b2 := c.b.Mark()
	c.b.Emit2(bytecode.Branch, 0)
	c.b.SetAddressHere(b1)

	if !c.t.IsNext("catch") {
		return c.NewTokenError("expected catch not found")
	}

	err = c.Statement()
	if err != nil {
		return err
	}
	c.b.SetAddressHere(b2)
	c.b.Emit1(bytecode.TryPop)

	return nil
}
