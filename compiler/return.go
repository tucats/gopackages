package compiler

import (
	"github.com/tucats/gopackages/bytecode"
)

// Return handles the return statment compilation
func (c *Compiler) Return() error {

	hasReturnValue := false
	if !c.StatementEnd() {
		bc, err := c.Expression()
		if err != nil {
			return err
		}
		if c.coerce.Mark() == 0 {
			return c.NewError(InvalidReturnValueError)
		}
		c.b.Append(bc)
		c.b.Append(c.coerce)
		hasReturnValue = true
	}

	// Stop execution of this stream
	c.b.Emit(bytecode.Return, hasReturnValue)
	return nil
}

// Exit handles the exit statment compilation
func (c *Compiler) Exit() error {

	c.b.Emit(bytecode.Load, "util")
	c.b.Emit(bytecode.Member, "Exit")

	argCount := 0
	if !c.StatementEnd() {
		bc, err := c.Expression()
		if err != nil {
			return err
		}
		c.b.Append(bc)
		argCount = 1
	}

	c.b.Emit(bytecode.Call, argCount)

	return nil
}
