package compiler

import (
	"github.com/tucats/gopackages/errors"
	"github.com/tucats/gopackages/expressions/bytecode"
	"github.com/tucats/gopackages/expressions/tokenizer"
)

// conditional handles parsing the ?: trinary operator. The first term is
// converted to a boolean value, and if true the second term is returned, else
// the third term. All terms must be present.
func (c *Compiler) conditional() error {
	// Parse the conditional
	err := c.logicalOr()
	if err != nil {
		return err
	}

	// If this is not a conditional, we're done. Conditionals
	// are only permitted when extensions are enabled.
	if c.t.AtEnd() || c.t.Peek(1) != tokenizer.OptionalToken {
		return nil
	}

	m1 := c.b.Mark()

	c.b.Emit(bytecode.BranchFalse, 0)

	// Parse both parts of the alternate values
	c.t.Advance(1)

	err = c.logicalOr()
	if err != nil {
		return err
	}

	if c.t.AtEnd() || c.t.Peek(1) != tokenizer.ColonToken {
		return c.error(errors.ErrMissingColon)
	}

	m2 := c.b.Mark()

	c.b.Emit(bytecode.Branch, 0)
	_ = c.b.SetAddressHere(m1)
	c.t.Advance(1)

	err = c.logicalOr()
	if err != nil {
		return err
	}

	// Patch up the forward references.
	_ = c.b.SetAddressHere(m2)

	return nil
}

func (c *Compiler) logicalAnd() error {
	err := c.relations()
	if err != nil {
		return err
	}

	for c.t.IsNext(tokenizer.BooleanAndToken) {
		// Handle short-circuit form boolean
		c.b.Emit(bytecode.Dup)

		mark := c.b.Mark()
		c.b.Emit(bytecode.BranchFalse, 0)

		err := c.relations()
		if err != nil {
			return err
		}

		c.b.Emit(bytecode.And)
		_ = c.b.SetAddressHere(mark)
	}

	return nil
}

func (c *Compiler) logicalOr() error {
	err := c.logicalAnd()
	if err != nil {
		return err
	}

	for c.t.IsNext(tokenizer.BooleanOrToken) {
		// Handle short-circuit from boolean
		c.b.Emit(bytecode.Dup)

		mark := c.b.Mark()
		c.b.Emit(bytecode.BranchTrue, 0)

		err := c.logicalAnd()
		if err != nil {
			return err
		}

		c.b.Emit(bytecode.Or)
		_ = c.b.SetAddressHere(mark)
	}

	return nil
}
