package compiler

import (
	"github.com/tucats/gopackages/errors"
	"github.com/tucats/gopackages/expressions/bytecode"
	"github.com/tucats/gopackages/expressions/data"
	"github.com/tucats/gopackages/expressions/tokenizer"
)

// reference parses a structure or array reference.
func (c *Compiler) reference() error {
	// Parse the function call or exprssion atom
	err := c.expressionAtom()
	if err != nil {
		return err
	}

	parsing := true
	// is there a trailing structure or array reference?
	for parsing && !c.t.AtEnd() {
		op := c.t.Peek(1)

		switch op {
		// Structure initialization
		case tokenizer.DataBeginToken:
			// If this is during switch statement processing, it can't be
			// a structure initialization.
			if c.flags.disallowStructInits {
				return nil
			}

			name := c.t.Peek(2)
			colon := c.t.Peek(3)

			if name.IsIdentifier() && colon == tokenizer.ColonToken {
				c.b.Emit(bytecode.Push, data.TypeMDKey)

				err := c.expressionAtom()
				if err != nil {
					return err
				}

				i := c.b.Opcodes()
				ix := i[len(i)-1]
				ix.Operand = data.Int(ix.Operand) + 1 // __type
				i[len(i)-1] = ix
			} else {
				parsing = false
			}
		// Function invocation
		case tokenizer.StartOfListToken:
			c.t.Advance(1)

			err := c.functionCall()
			if err != nil {
				return err
			}

		// Array index reference
		case tokenizer.StartOfArrayToken:
			c.t.Advance(1)

			// If there is an slice with an implied start of 0,
			// handle that here.
			t := c.t.Peek(1)
			if t == tokenizer.ColonToken {
				c.b.Emit(bytecode.Push, 0)
			} else {
				err := c.conditional()
				if err != nil {
					return err
				}
			}

			// is it a slice instead of an index?
			if c.t.IsNext(tokenizer.ColonToken) {
				// IS this the case of the assumed end being the
				// length of the item? If so, add code to use the
				// length of the item below current ToS. The actual
				// displacement is 2, since before executing it we
				// also already pushed the length fuction on stack.
				if c.t.Peek(1) == tokenizer.EndOfArrayToken {
					c.b.Emit(bytecode.Load, "len")
					c.b.Emit(bytecode.ReadStack, -2)
					c.b.Emit(bytecode.Call, 1)
				} else {
					err := c.conditional()
					if err != nil {
						return err
					}
				}

				c.b.Emit(bytecode.LoadSlice)

				if c.t.Next() != tokenizer.EndOfArrayToken {
					return c.error(errors.ErrMissingBracket)
				}
			} else {
				// Nope, singular index
				if c.t.Next() != tokenizer.EndOfArrayToken {
					return c.error(errors.ErrMissingBracket)
				}

				c.b.Emit(bytecode.LoadIndex)
			}

		// Nothing else, term is complete
		default:
			return nil
		}
	}

	return nil
}
