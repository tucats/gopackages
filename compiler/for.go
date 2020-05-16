package compiler

import (
	"errors"

	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/expressions"
)

// For compiles the loop statement. This has three clauses
// which are separated by ";", followed by a statement or
// block that is run as described by the loop conditions.
func (c *Compiler) For() error {

	// Crude assignment statement test
	if !c.IsLValue() {
		return c.NewError("loop initialization not found")
	}

	lv, err := c.LValue()
	if err != nil {
		return err
	}
	if !c.t.IsNext(":=") {
		return errors.New("expected := not found")
	}

	bc, err := expressions.Compile(c.t)
	if err != nil {
		return err
	}
	c.b.Append(bc)
	c.b.Append(lv)

	if !c.t.IsNext(";") {
		c.NewError("missing ; in loop definition")
	}

	condition, err := expressions.Compile(c.t)
	if err != nil {
		return err
	}

	if !c.t.IsNext(";") {
		c.NewError("missing ; in loop definition")
	}

	lv, err = c.LValue()
	if err != nil {
		return err
	}

	if !c.t.IsNext(":=") {
		return errors.New("expected := not found")
	}

	bc, err = expressions.Compile(c.t)
	if err != nil {
		return err
	}

	// Top of loop starts here
	b1 := c.b.Mark()

	// Test condition
	c.b.Append(condition)
	b2 := c.b.Mark()
	c.b.Emit2(bytecode.BranchFalse, 0)

	// Loop body goes next
	err = c.Statement()
	if err != nil {
		return err
	}

	// Emit increment code, and loop. Finally, mark the exit location from
	// the condition test for the loop.
	c.b.Append(bc)
	c.b.Append(lv)
	c.b.Emit2(bytecode.Branch, b1)
	c.b.SetAddressHere(b2)

	return nil
}
