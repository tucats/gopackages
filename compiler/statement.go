package compiler

import (
	"errors"

	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/expressions"
)

// Statement parses a single statement
func (c *Compiler) Statement() error {

	// Statement block
	if c.t.IsNext("{") {
		return c.Block()
	}

	// Crude assignment statement test
	if c.t.Peek(2) == ":=" {
		name := c.t.Next()
		c.t.Advance(1)
		bc, err := expressions.Compile(c.t)
		if err != nil {
			return err
		}
		c.b.Append(bc)
		c.b.Emit(bytecode.Store, name)
		return nil
	}

	if c.t.IsNext("print") {
		return c.Print()
	}
	
	return errors.New("unrecognized statement")
}
