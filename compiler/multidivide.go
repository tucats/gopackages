package compiler

import (
	bc "github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/tokenizer"
)

// Eval evaluates the parsed expression. This can be called multiple times
// with the same scanned string, but with different symbols.
func (c *Compiler) multDivide() error {

	err := c.unary()
	if err != nil {
		return err
	}

	var parsing = true
	for parsing {
		if c.t.AtEnd() {
			break
		}
		op := c.t.Peek(1)
		if c.t.AnyNext([]string{"^", "*", "/", "|"}) {

			if c.t.IsNext(tokenizer.EndOfTokens) {
				return c.NewError(MissingTermError)
			}

			err := c.unary()
			if err != nil {
				return err
			}

			switch op {

			case "^":
				c.b.Emit(bc.Exp)

			case "*":
				c.b.Emit(bc.Mul)

			case "/":
				c.b.Emit(bc.Div)

			case "|":
				c.b.Emit(bc.Or)

			}

		} else {
			parsing = false
		}
	}
	return nil
}
