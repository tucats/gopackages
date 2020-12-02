package compiler

import (
	bc "github.com/tucats/gopackages/bytecode"
)

func (c *Compiler) relations() error {

	err := c.addSubtract()
	if err != nil {
		return err
	}

	var parsing = true
	for parsing {
		if c.t.AtEnd() {
			break
		}
		op := c.t.Peek(1)
		if op == "==" || op == "!=" || op == "<" || op == "<=" || op == ">" || op == ">=" {
			c.t.Advance(1)

			err := c.addSubtract()
			if err != nil {
				return err
			}

			switch op {

			case "==":
				c.b.Emit(bc.Equal)

			case "!=":
				c.b.Emit(bc.NotEqual)

			case "<":
				c.b.Emit(bc.LessThan)

			case "<=":
				c.b.Emit(bc.LessThanOrEqual)

			case ">":
				c.b.Emit(bc.GreaterThan)

			case ">=":
				c.b.Emit(bc.GreaterThanOrEqual)

			}

		} else {
			parsing = false
		}
	}
	return nil
}
