package compiler

import (
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/expressions"
)

// Return handles the return statment compilation
func (c *Compiler) Return() error {

	if !c.StatementEnd() {
		bc, err := expressions.Compile(c.t)
		if err != nil {
			return err
		}
		if c.coerce.Mark() == 0 {
			return c.NewTokenError("return value from void function")
		}
		c.b.Append(bc)
	}

	// Is there a coerce to set to the required type?
	c.b.Append(c.coerce)

	// Stop execution of this stream
	c.b.Emit1(bytecode.Stop)

	return nil
}

// Exit handles the exit statment compilation
func (c *Compiler) Exit() error {

	c.b.Emit2(bytecode.Load, "util")
	c.b.Emit2(bytecode.Member, "exit")

	argCount := 0
	if !c.StatementEnd() {
		bc, err := expressions.Compile(c.t)
		if err != nil {
			return err
		}
		c.b.Append(bc)
		argCount = 1
	}

	c.b.Emit2(bytecode.Call, argCount)

	return nil
}
