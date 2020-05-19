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

	indexName := ""
	// Is this the two-value range thing?
	if tokenizer.IsSymbol(c.t.Peek(1)) && (c.t.Peek(2) == ",") {
		indexName = c.t.Next()
		c.t.Advance(1)
	}
	// Must be a valid lvalue
	if !c.IsLValue() {
		return c.NewError("loop initialization not found")
	}

	indexStore, err := c.LValue()
	if err != nil {
		return err
	}
	if !c.t.IsNext(":=") {
		return errors.New("expected := not found")
	}

	// Do we compile a range?
	if c.t.IsNext("range") {

		c.PushLoop(rangeLoopType)

		arrayCode, err := expressions.Compile(c.t)
		if err != nil {
			return err
		}

		// Initialize index
		if indexName == "" {
			indexName = MakeSymbol()
		}

		c.b.Emit2(bytecode.Push, 1)
		c.b.Emit2(bytecode.Store, indexName)

		// Remember top of loop
		b1 := c.b.Mark()

		// Is index >= len of array?
		c.b.Emit2(bytecode.Load, "len")
		c.b.Append(arrayCode)
		c.b.Emit2(bytecode.Call, 1)
		c.b.Emit2(bytecode.Load, indexName)
		c.b.Emit1(bytecode.LessThan)

		b2 := c.b.Mark()
		c.b.Emit2(bytecode.BranchTrue, 0)

		// Load element of array
		c.b.Append(arrayCode)
		c.b.Emit2(bytecode.Load, indexName)
		c.b.Emit1(bytecode.LoadIndex)
		c.b.Append(indexStore)

		err = c.Statement()
		if err != nil {
			return err
		}

		// Increment the index
		b3 := c.b.Mark()
		c.b.Emit2(bytecode.Load, indexName)
		c.b.Emit2(bytecode.Push, 1)
		c.b.Emit1(bytecode.Add)
		c.b.Emit2(bytecode.Store, indexName)

		// Branch back to start of loop
		c.b.Emit2(bytecode.Branch, b1)
		for _, fixAddr := range c.loops.continues {
			c.b.SetAddress(fixAddr, b3)
		}

		c.b.SetAddressHere(b2)

		for _, fixAddr := range c.loops.breaks {
			c.b.SetAddressHere(fixAddr)
		}
		c.PopLoop()
		c.b.Emit2(bytecode.SymbolDelete, indexName)
		return nil
	}

	// Nope, normal numeric loop ocnditions. IF so, it cannot
	// have an index variable defined.
	if indexName != "" {
		c.NewError("invalid index variable")
	}
	c.PushLoop(indexLoopType)

	// The expression is the initial value of the loop.
	initializerCode, err := expressions.Compile(c.t)
	if err != nil {
		return err
	}
	c.b.Append(initializerCode)
	c.b.Append(indexStore)

	if !c.t.IsNext(";") {
		c.NewError("missing ; in loop definition")
	}

	// Now get the condition clause that tells us if the loop
	// is still executing.
	condition, err := expressions.Compile(c.t)
	if err != nil {
		return err
	}

	if !c.t.IsNext(";") {
		c.NewError("missing ; in loop definition")
	}

	// Finally, get the clause that updates something
	// (nominally the index) to eventuall trigger the
	// loop condition.
	incrementStore, err := c.LValue()
	if err != nil {
		return err
	}

	if !c.t.IsNext(":=") {
		return errors.New("expected := not found")
	}

	incrementCode, err := expressions.Compile(c.t)
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
	c.b.Append(incrementCode)
	c.b.Append(incrementStore)
	c.b.Emit2(bytecode.Branch, b1)
	c.b.SetAddressHere(b2)

	for _, fixAddr := range c.loops.continues {
		c.b.SetAddress(fixAddr, b1)
	}

	for _, fixAddr := range c.loops.breaks {
		c.b.SetAddressHere(fixAddr)
	}

	return nil
}

// Break processes a break statement
func (c *Compiler) Break() error {

	if c.loops == nil {
		return c.NewError("break outside of loop")
	}
	fixAddr := c.b.Mark()
	c.b.Emit2(bytecode.Branch, 0)
	c.loops.breaks = append(c.loops.breaks, fixAddr)
	return nil
}

// Continue processes a continue statement
func (c *Compiler) Continue() error {

	if c.loops == nil {
		return c.NewError("continue outside of loop")
	}
	fixAddr := c.b.Mark()
	c.b.Emit2(bytecode.Branch, 0)
	c.loops.continues = append(c.loops.continues, fixAddr)
	return nil
}

// PushLoop creates a new loop context and adds it to the
// top of the loop stack.
func (c *Compiler) PushLoop(loopType int) {

	loop := Loop{
		Type:      loopType,
		breaks:    make([]int, 0),
		continues: make([]int, 0),
		Parent:    c.loops,
	}

	c.loops = &loop
}

// PopLoop discards the topmost loop on the loop stack.
func (c *Compiler) PopLoop() {
	if c.loops != nil {
		c.loops = c.loops.Parent
	}
}
