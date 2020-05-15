package compiler

import (
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/expressions"
)

// Call handles the call statement.
func (c *Compiler) Call() error {

	bc, err := expressions.Compile(c.t)
	if err != nil {
		return err
	}
	c.b.Append(bc)
	c.b.Emit0(bytecode.Drop)
	return nil
}
