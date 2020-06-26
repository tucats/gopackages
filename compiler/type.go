package compiler

import (
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/tokenizer"
)

// Type compiles a type statement
func (c *Compiler) Type() error {

	name := c.t.Next()
	if !tokenizer.IsSymbol(name) {
		return c.NewError("invalid type name")
	}

	parent := name
	if c.t.Peek(1) == "->" {
		c.t.Advance(1)
		parent = c.t.Next()
		if !tokenizer.IsSymbol(parent) {
			return c.NewError("invalid parent name")
		}
	}
	if c.t.Peek(1) != "{" {
		return c.NewTokenError("expected {, found ")
	}

	// If there is not parent, seal the chain by making the link point to a string of our own name.
	// If there is a parent, load it so it can be linked after type creation.
	if parent == name {
		c.b.Emit2(bytecode.Push, parent)
	} else {
		c.b.Emit2(bytecode.Load, parent)
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
	c.b.Emit2(bytecode.Push, "__parent")
	c.b.Emit1(bytecode.StoreIndex)
	c.b.Emit2(bytecode.SymbolCreate, name)
	c.b.Emit2(bytecode.Store, name)

	return nil
}

func (c *Compiler) compileType() error {

	// Must start with {
	if !c.t.IsNext("{") {
		return c.NewError("expected { not found")
	}

	count := 0
	for true {
		name := c.t.Next()
		if !tokenizer.IsSymbol(name) {
			return c.NewStringError("invalid member name", name)
		}
		count = count + 1
		if c.t.Peek(1) == "{" {
			err := c.compileType()
			if err != nil {
				return err
			}
		} else {
			switch c.t.Next() {
			case "int":
				c.b.Emit2(bytecode.Push, 0)
			case "float":
				c.b.Emit2(bytecode.Push, 0.0)
			case "bool":
				c.b.Emit2(bytecode.Push, false)
			case "string":
				c.b.Emit2(bytecode.Push, "")
			default:
				return c.NewTokenError("invalid type")
			}
		}

		c.b.Emit2(bytecode.Push, name)

		if c.t.IsNext("}") {
			c.b.Emit2(bytecode.Struct, count)
			return nil
		}
		if c.t.AtEnd() {
			return c.NewError("incomplete type definition")
		}
	}
	return nil
}
