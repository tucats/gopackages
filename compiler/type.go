package compiler

import (
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/tokenizer"
)

// Type compiles a type statement
func (c *Compiler) Type() error {

	name := c.t.Next()
	if !tokenizer.IsSymbol(name) {
		return c.NewError(InvalidSymbolError)
	}
	name = c.Normalize(name)

	parent := name
	if c.t.Peek(1) == "->" {
		c.t.Advance(1)
		parent = c.t.Next()
		if !tokenizer.IsSymbol(parent) {
			return c.NewError(InvalidSymbolError)
		}
		c.Normalize(parent)
	}

	// Make sure this is a legit type definition
	if c.t.Peek(1) == "struct" && c.t.Peek(2) == "{" {
		c.t.Advance(1)
	}
	if c.t.Peek(1) != "{" {
		return c.NewError(MissingBlockError)
	}

	// If there is not parent, seal the chain by making the link point to a string of our own name.
	// If there is a parent, load it so it can be linked after type creation.
	if parent == name {
		c.b.Emit(bytecode.Push, parent)
	} else {
		c.b.Emit(bytecode.Load, parent)
	}

	// Compile a struct definition
	err := c.compileType()
	if err != nil {
		return err
	}

	// Add in the type linkage, and store as the type name. The __parent for a type is
	// a string that is the name of the type. When a member dereference on a struct
	// happens that includes a __parent, the __parent object is also checked for the
	// member if it is NOT a string.
	c.b.Emit(bytecode.Push, "__parent")
	c.b.Emit(bytecode.StoreIndex, true)
	c.b.Emit(bytecode.SymbolCreate, name)
	c.b.Emit(bytecode.Dup)
	c.b.Emit(bytecode.Dup)
	c.b.Emit(bytecode.Store, name)

	// Use the name as the type.
	c.b.Emit(bytecode.Push, name)
	c.b.Emit(bytecode.Swap)
	c.b.Emit(bytecode.Push, "__type")
	c.b.Emit(bytecode.StoreIndex, true)

	// Finally, make it a static value now.
	c.b.Emit(bytecode.Dup) // One more needed for type statement
	c.b.Emit(bytecode.Push, true)
	c.b.Emit(bytecode.Swap)
	c.b.Emit(bytecode.Push, "__static")
	c.b.Emit(bytecode.StoreIndex, true)

	return nil
}

func (c *Compiler) compileType() error {

	// Skip over the optional struct type keyword
	if c.t.Peek(1) == "struct" && c.t.Peek(2) == "{" {
		c.t.Advance(1)
	}

	// Must start with {
	if !c.t.IsNext("{") {
		return c.NewError(MissingBlockError)
	}

	count := 0
	for {
		name := c.t.Next()
		if !tokenizer.IsSymbol(name) {
			return c.NewError(InvalidSymbolError, name)
		}
		name = c.Normalize(name)

		count = count + 1
		// Skip over the optional struct type keyword
		if c.t.Peek(1) == "struct" && c.t.Peek(2) == "{" {
			c.t.Advance(1)
		}
		if c.t.Peek(1) == "{" {
			err := c.compileType()
			if err != nil {
				return err
			}
		} else {
			switch c.t.Next() {
			case "int":
				c.b.Emit(bytecode.Push, 0)
			case "float":
				c.b.Emit(bytecode.Push, 0.0)
			case "bool":
				c.b.Emit(bytecode.Push, false)
			case "string":
				c.b.Emit(bytecode.Push, "")
			default:
				return c.NewError(InvalidTypeNameError)
			}
		}

		c.b.Emit(bytecode.Push, name)

		// Eat any trailing commas, and the see if we're at the end
		_ = c.t.IsNext(",")
		if c.t.IsNext("}") {
			c.b.Emit(bytecode.Struct, count)
			return nil
		}
		if c.t.AtEnd() {
			return c.NewError(MissingEndOfBlockError)
		}
	}
}
