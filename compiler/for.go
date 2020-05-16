package compiler

import (
	"errors"

	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/expressions"
	"github.com/tucats/gopackages/tokenizer"
)

// For compiles the loop statement. This has three clauses
// which are separated by ";", followed by a statement or
// block that is run as described by the loop conditions.
func (c *Compiler) For() error {

	index := ""
	// Is this the two-value range thing?
	if tokenizer.IsSymbol(c.t.Peek(1)) && (c.t.Peek(2) == ",") {
		index = c.t.Next()
		c.t.Advance(1)
	}
	// Must be a valid lvalue
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

	// Here is where I should check for "range"
	//
	// for n := range array {
	//	...
	// }
	//
	// for _i := 0; _i < len(array); _i = _i + 1 {
	//     n := array[_i]
	//	   ...
	// }

	// Do we compile a range?
	if c.t.IsNext("range") {

		ex, err := expressions.Compile(c.t)
		if err != nil {
			return err
		}

		// Initialize index
		if index == "" {
			index = MakeSymbol()
		}
		c.b.Emit2(bytecode.Push, 1)
		c.b.Emit2(bytecode.Store, index)

		// Remember top of loop
		b1 := c.b.Mark()

		// Is index >= len of array?
		c.b.Append(ex)
		c.b.Emit2(bytecode.Push, "len")
		c.b.Emit2(bytecode.Call, 1)
		c.b.Emit2(bytecode.Load, index)
		c.b.Emit1(bytecode.LessThan)

		b2 := c.b.Mark()
		c.b.Emit2(bytecode.BranchTrue, 0)

		// Load element of array
		c.b.Append(ex)
		c.b.Emit2(bytecode.Load, index)
		c.b.Emit1(bytecode.LoadIndex)
		c.b.Append(lv)

		err = c.Statement()
		if err != nil {
			return err
		}

		// Increment the index
		c.b.Emit2(bytecode.Load, index)
		c.b.Emit2(bytecode.Push, 1)
		c.b.Emit1(bytecode.Add)
		c.b.Emit2(bytecode.Store, index)

		// Branch back to start of loop
		c.b.Emit2(bytecode.Branch, b1)
		c.b.SetAddressHere(b2)
		return nil
	}

	// Nope, normal numeric loop ocnditions. IF so, it cannot
	// have an index variable defined.
	if index != "" {
		c.NewError("invalid index variable")
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
