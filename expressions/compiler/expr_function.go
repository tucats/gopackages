package compiler

import (
	"github.com/tucats/gopackages/errors"
	bc "github.com/tucats/gopackages/expressions/bytecode"
	"github.com/tucats/gopackages/expressions/tokenizer"
)

func (c *Compiler) functionCall() error {
	// Note, caller already consumed the opening paren
	argc := 0

	for c.t.Peek(1) != tokenizer.EndOfListToken {
		err := c.conditional()
		if err != nil {
			return err
		}

		argc = argc + 1

		if c.t.AtEnd() {
			break
		}

		if c.t.Peek(1) == tokenizer.EndOfListToken {
			break
		}

		// Could be the "..." flatten operator
		if c.t.IsNext(tokenizer.VariadicToken) {
			c.b.Emit(bc.Flatten)

			break
		}

		if c.t.Peek(1) != tokenizer.CommaToken {
			return c.error(errors.ErrInvalidList)
		}

		c.t.Advance(1)
	}

	// Ensure trailing parenthesis
	if c.t.AtEnd() || c.t.Peek(1) != tokenizer.EndOfListToken {
		return c.error(errors.ErrMissingParenthesis)
	}

	c.t.Advance(1)
	// Call the function
	c.b.Emit(bc.Call, argc)

	return nil
}

// functionOrReference compiles a function call. The value of the
// function has been pushed to the top of the stack.
func (c *Compiler) functionOrReference() error {
	// Get the atom
	err := c.reference()
	if err != nil {
		return err
	}

	// Peek ahead to see if it's the start of a function call...
	if c.t.IsNext(tokenizer.StartOfListToken) {
		return c.functionCall()
	}

	return nil
}
