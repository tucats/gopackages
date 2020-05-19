package compiler

import (
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/expressions"
	"github.com/tucats/gopackages/tokenizer"
)

// Array compiles the array statement
func (c *Compiler) Array() error {

	name := c.t.Next()
	if !tokenizer.IsSymbol(name) {
		c.t.Advance(-1)
		return c.NewTokenError("invalid array name")
	}
	// See  if it's on a reserved word.
	if tokenizer.InList(name, []string{"print", "for", "array", "if", "call", "return"}) {
		c.t.Advance(-1)
		return c.NewTokenError("invalid array name")
	}

	if !c.t.IsNext("[") {
		return c.NewError("missing [ in array")
	}

	bc, err := expressions.Compile(c.t)
	if err != nil {
		return err
	}
	c.b.Append(bc)
	if !c.t.IsNext("]") {
		return c.NewError("missing ] in array")
	}
	if c.t.IsNext("=") {
		bc, err = expressions.Compile(c.t)
		if err != nil {
			return nil
		}
		c.b.Append(bc)
		c.b.Emit2(bytecode.MakeArray, 2)
	} else {
		c.b.Emit2(bytecode.MakeArray, 1)
	}
	c.b.Emit2(bytecode.Store, name)
	return nil
}
