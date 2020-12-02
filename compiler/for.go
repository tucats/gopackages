package compiler

import (
	"errors"

	"github.com/tucats/gopackages/app-cli/ui"
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/tokenizer"
)

// For compiles the loop statement. This has two syntax types that
// can be specified.
// 1. There are three clauses which are separated by ";", followed
//    by a statement or block that is run as described by the loop
//    index variable conditions.
//
// 2. There can be a range operation which creates an implied loop
//    using each member of the array or struct.
func (c *Compiler) For() error {

	c.b.Emit(bytecode.PushScope)

	// Is this the two-value range thing?
	indexName := ""
	if tokenizer.IsSymbol(c.t.Peek(1)) && (c.t.Peek(2) == ",") {
		indexName = c.t.Next()
		c.t.Advance(1)
	}

	// Must be a valid lvalue
	if !c.IsLValue() {
		return c.NewError(MissingForLoopInitializerError)
	}

	indexStore, err := c.LValue()
	if err != nil {
		return err
	}
	if !c.t.IsNext(":=") {
		return errors.New(MissingLoopAssignmentError)
	}

	// Do we compile a range?
	if c.t.IsNext("range") {

		// This is wierd, but the LValue compiler will have inserted a "SymbolCreate" in the
		// lValue due to the syntax, but we don't really want to create it as it will have already
		// been generated once. So use it once to create a value, and then remove the store.
		c.b.Emit(bytecode.Push, 0)
		c.b.Append(indexStore)
		indexStore.Remove(0)

		// Make a new scope and get the array we will range over
		c.PushLoop(rangeLoopType)
		arrayCode, err := c.Expression()
		if err != nil {
			return err
		}

		// Initialize index
		if indexName == "" {
			indexName = MakeSymbol()
		}
		c.b.Emit(bytecode.Push, 0)
		c.b.Emit(bytecode.SymbolCreate, indexName)
		c.b.Emit(bytecode.Store, indexName)

		// Remember top of loop
		b1 := c.b.Mark()

		// Is index >= len of array?
		c.b.Emit(bytecode.Load, "len")
		c.b.Append(arrayCode)
		c.b.Emit(bytecode.Call, 1)
		c.b.Emit(bytecode.Load, indexName)
		c.b.Emit(bytecode.LessThanOrEqual)

		b2 := c.b.Mark()
		c.b.Emit(bytecode.BranchTrue, 0)

		// Load element of array
		c.b.Append(arrayCode)
		c.b.Emit(bytecode.Load, indexName)
		c.b.Emit(bytecode.LoadIndex)
		c.b.Append(indexStore)

		err = c.Statement()
		if err != nil {
			return err
		}

		// Increment the index
		b3 := c.b.Mark()
		c.b.Emit(bytecode.Load, indexName)
		c.b.Emit(bytecode.Push, 1)
		c.b.Emit(bytecode.Add)
		c.b.Emit(bytecode.Store, indexName)

		// Branch back to start of loop
		c.b.Emit(bytecode.Branch, b1)
		for _, fixAddr := range c.loops.continues {
			c.b.SetAddress(fixAddr, b3)
		}

		c.b.SetAddressHere(b2)

		for _, fixAddr := range c.loops.breaks {
			c.b.SetAddressHere(fixAddr)
		}
		c.PopLoop()
		c.b.Emit(bytecode.SymbolDelete, indexName)
		c.b.Emit(bytecode.PopScope)

		return nil
	}

	// Nope, normal numeric loop conditions. At this point there should not
	// be an index variable defined.
	if indexName != "" {
		c.NewError(InvalidLoopIndexError)
	}
	c.PushLoop(indexLoopType)

	// The expression is the initial value of the loop.
	initializerCode, err := c.Expression()
	if err != nil {
		return err
	}
	c.b.Append(initializerCode)
	c.b.Append(indexStore)

	if !c.t.IsNext(";") {
		c.NewError(MissingSemicolonError)
	}

	// Now get the condition clause that tells us if the loop
	// is still executing.
	condition, err := c.Expression()
	if err != nil {
		return err
	}

	if !c.t.IsNext(";") {
		c.NewError(MissingSemicolonError)
	}

	// Finally, get the clause that updates something
	// (nominally the index) to eventually trigger the
	// loop condition.
	incrementStore, err := c.LValue()
	if err != nil {
		return err
	}

	if !c.t.IsNext("=") {
		return errors.New(MissingEqualError)
	}

	incrementCode, err := c.Expression()
	if err != nil {
		return err
	}

	// Top of loop body starts here
	b1 := c.b.Mark()

	// Emit the test condition
	c.b.Append(condition)
	b2 := c.b.Mark()
	c.b.Emit(bytecode.BranchFalse, 0)

	// Loop body goes next
	err = c.Statement()
	if err != nil {
		return err
	}

	// Emit increment code, and loop. Finally, mark the exit location from
	// the condition test for the loop.
	c.b.Append(incrementCode)
	c.b.Append(incrementStore)
	c.b.Emit(bytecode.Branch, b1)
	c.b.SetAddressHere(b2)

	for _, fixAddr := range c.loops.continues {
		c.b.SetAddress(fixAddr, b1)
	}

	for _, fixAddr := range c.loops.breaks {
		c.b.SetAddressHere(fixAddr)
	}
	c.b.Emit(bytecode.PopScope)
	c.PopLoop()
	return nil
}

// Break compiles a break statement
func (c *Compiler) Break() error {
	if c.loops == nil {
		return c.NewError(InvalidLoopControlError)
	}
	fixAddr := c.b.Mark()
	c.b.Emit(bytecode.Branch, 0)
	c.loops.breaks = append(c.loops.breaks, fixAddr)
	return nil
}

// Continue compiles a continue statement
func (c *Compiler) Continue() error {
	if c.loops == nil {
		return c.NewError(InvalidLoopControlError)
	}
	fixAddr := c.b.Mark()
	c.b.Emit(bytecode.Branch, 0)
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
	// ui.Debug("=== Push loop scope")
	c.loops = &loop
}

// PopLoop discards the top-most loop on the loop stack.
func (c *Compiler) PopLoop() {
	if c.loops != nil {
		c.loops = c.loops.Parent
		// ui.Debug("=== Pop loop scope")
	} else {
		ui.Debug("=== loop stack empty")
	}
}
