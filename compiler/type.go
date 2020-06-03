package compiler

import (
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/expressions"
	"github.com/tucats/gopackages/tokenizer"
)

// Type compiles a type statement
func (c *Compiler) Type() error {

	name := c.t.Next()
	if !tokenizer.IsSymbol(name) {
		return c.NewError("invalid type name")
	}

	if c.t.Peek(1) != "{" {
		return c.NewTokenError("expected {, found ")
	}

	// Compile a struct definition
	ex, err := expressions.Compile(c.t)
	if err != nil {
		return err
	}

	c.b.Emit2(bytecode.Push, name)
	c.b.Append(ex)

	// Add in the type linkage, and store as the type name. The __type for a type is
	// a string that is the name of the type. When a member dereference on a struct
	// happens that includes a __type, the __type object is also checked for the
	// member if it is NOT a string.
	c.b.Emit2(bytecode.Push, "__type")
	c.b.Emit1(bytecode.StoreIndex)
	c.b.Emit2(bytecode.SymbolCreate, name)
	c.b.Emit2(bytecode.Store, name)

	return nil
}
