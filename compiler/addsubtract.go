package compiler

import (
	bc "github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/tokenizer"
)

func (c *Compiler) addSubtract() error {

	err := c.multDivide()
	if err != nil {
		return err
	}

	var parsing = true
	for parsing {
		if c.t.AtEnd() {
			break
		}
		op := c.t.Peek(1)
		if tokenizer.InList(op, []string{"+", "-", "&"}) {
			c.t.Advance(1)

			if c.t.IsNext(tokenizer.EndOfTokens) {
				return c.NewError(MissingTermError)
			}

			err := c.multDivide()
			if err != nil {
				return err
			}

			switch op {

			case "+":
				c.b.Emit(bc.Add)

			case "-":
				c.b.Emit(bc.Sub)

			case "&":
				c.b.Emit(bc.And)
			}

		} else {
			parsing = false
		}
	}
	return nil
}
