package compiler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tucats/gopackages/defs"
	"github.com/tucats/gopackages/errors"
	"github.com/tucats/gopackages/expressions/bytecode"
	"github.com/tucats/gopackages/expressions/tokenizer"
)

func (c *Compiler) expressionAtom() error {
	t := c.t.Peek(1)

	// Is it a binary constant? If so, convert to decimal.
	text := t.Spelling()
	if strings.HasPrefix(strings.ToLower(text), "0b") {
		binaryValue := 0
		fmt.Sscanf(text[2:], "%b", &binaryValue)
		t = tokenizer.NewIntegerToken(strconv.Itoa(binaryValue))
	}

	// Is it a hexadecimal constant? If so, convert to decimal.
	if strings.HasPrefix(strings.ToLower(text), "0x") {
		hexValue := 0
		fmt.Sscanf(strings.ToLower(text[2:]), "%x", &hexValue)
		t = tokenizer.NewIntegerToken(strconv.Itoa(hexValue))
	}

	// Is it an octal constant? If so, convert to decimal.
	if strings.HasPrefix(strings.ToLower(text), "0o") {
		octalValue := 0
		fmt.Sscanf(strings.ToLower(text[2:]), "%o", &octalValue)
		t = tokenizer.NewIntegerToken(strconv.Itoa(octalValue))
	}

	// Is this the "nil" constant?
	if t == tokenizer.NilToken {
		c.t.Advance(1)
		c.b.Emit(bytecode.Push, nil)

		return nil
	}

	// Is this a parenthesis expression?
	if t == tokenizer.StartOfListToken {
		c.t.Advance(1)

		err := c.conditional()
		if err != nil {
			return err
		}

		if c.t.Next() != tokenizer.EndOfListToken {
			return c.error(errors.ErrMissingParenthesis)
		}

		return nil
	}

	// If the token is a number, convert it
	if t.IsClass(tokenizer.IntegerTokenClass) {
		if i, err := strconv.ParseInt(text, 10, 32); err == nil {
			c.t.Advance(1)
			c.b.Emit(bytecode.Push, int(i))

			return nil
		}

		if i, err := strconv.ParseInt(text, 10, 64); err == nil {
			c.t.Advance(1)
			c.b.Emit(bytecode.Push, i)

			return nil
		}
	}

	if t.IsClass(tokenizer.FloatTokenClass) {
		if i, err := strconv.ParseFloat(text, 64); err == nil {
			c.t.Advance(1)
			c.b.Emit(bytecode.Push, i)

			return nil
		}
	}

	if t.IsClass(tokenizer.BooleanTokenClass) {
		if text == defs.True || text == defs.False {
			c.t.Advance(1)
			c.b.Emit(bytecode.Push, (text == defs.True))

			return nil
		}
	}

	if t.IsValue() {
		c.t.Advance(1)
		c.b.Emit(bytecode.Push, t)

		return nil
	}

	// Is it just a symbol needing a load?
	if t.IsIdentifier() {
		c.b.Emit(bytecode.Load, t)
		c.t.Advance(1)

		return nil
	}

	// Not something we know what to do with...
	return c.error(errors.ErrUnexpectedToken, t)
}
