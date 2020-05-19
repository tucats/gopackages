package compiler

import "github.com/tucats/gopackages/bytecode"

// Block compiles a statement block. The leading { has already
// been parse.
func (c *Compiler) Block() error {

	parsing := true
	c.b.Emit1(bytecode.PushScope)
	for parsing {

		if c.t.IsNext("}") {
			break
		}

		err := c.Statement()
		if err != nil {
			return err
		}

		if c.t.IsNext(";") {
			// No action needed
		}

		if c.t.AtEnd() {
			return c.NewError("unclosed statement block")
		}
	}
	c.b.Emit1(bytecode.PopScope)
	return nil
}
