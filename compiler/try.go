package compiler

import "github.com/tucats/gopackages/bytecode"

// Try compiles the try statement which allows the program to catch error
// conditions instead of stopping execution on an error.
func (c *Compiler) Try() error {

	// Generate start of a try block.
	b1 := c.b.Mark()
	c.b.Emit(bytecode.Try, 0)

	// Statement to try
	err := c.Statement()
	if err != nil {
		return err
	}
	b2 := c.b.Mark()
	c.b.Emit(bytecode.Branch, 0)
	_ = c.b.SetAddressHere(b1)

	if !c.t.IsNext("catch") {
		return c.NewError(MissingCatchError)
	}

	err = c.Statement()
	if err != nil {
		return err
	}
	_ = c.b.SetAddressHere(b2)
	c.b.Emit(bytecode.TryPop)

	return nil
}
