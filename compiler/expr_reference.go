package compiler

import (
	bc "github.com/tucats/gopackages/bytecode"
)

// reference parses a structure or array reference
func (c *Compiler) reference() error {

	// Parse the function call or exprssion atom
	err := c.expressionAtom()
	if err != nil {
		return err
	}

	// is there a trailing structure or array reference?
	for !c.t.AtEnd() {

		op := c.t.Peek(1)
		switch op {

		// Struct reference
		case "->":
			c.t.Advance(1)
			name := c.t.Next()
			c.b.Emit(bc.Dup)
			c.b.Emit(bc.Push, name)
			c.b.Emit(bc.ClassMember)

		// Map member reference
		case ".":
			c.t.Advance(1)
			name := c.t.Next()
			c.b.Emit(bc.Push, name)
			c.b.Emit(bc.Member)

		// Array index reference
		case "[":
			c.t.Advance(1)
			err := c.conditional()
			if err != nil {
				return err
			}

			// is it a slice instead of an index?
			if c.t.IsNext(":") {
				err := c.conditional()
				if err != nil {
					return err
				}
				c.b.Emit(bc.LoadSlice)
				if c.t.Next() != "]" {
					return c.NewError(MissingBracketError)
				}
			} else {
				// Nope, singular index
				if c.t.Next() != "]" {
					return c.NewError(MissingBracketError)
				}
				c.b.Emit(bc.LoadIndex)
			}

		// Nothing else, term is complete
		default:
			return nil
		}
	}
	return nil
}
