package compiler

import (
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/expressions"
	"github.com/tucats/gopackages/tokenizer"
)

// IsLValue peeks ahead to see if this is likely to be an lValue
// object. This is used in cases where the parser might be in an
// otherwise ambiguous state
func (c *Compiler) IsLValue() bool {
	name := c.t.Peek(1)
	if !tokenizer.IsSymbol(name) {
		return false
	}

	next := c.t.Peek(2)
	if next == "." || next == "[" {
		return true
	}

	if next == ":=" {
		return true
	}
	return false
}

// LValue compiles the informaiton on the left side of
// an assignment. This information is used later to store the
// data in the named object.
func (c *Compiler) LValue() (*bytecode.ByteCode, error) {

	bc := bytecode.New("lvalue")
	name := c.t.Next()

	if !tokenizer.IsSymbol(name) {
		return nil, c.NewTokenError("invalid symbol name")
	}

	needLoad := true
	// Until we get to the end of the lvalue...
	for c.t.Peek(1) == "." || c.t.Peek(1) == "[" {

		if needLoad {
			bc.Emit(bytecode.Load, name)
			needLoad = false
		}
		err := c.lvalueTerm(bc)
		if err != nil {
			return nil, err
		}

	}

	bc.Emit(bytecode.Store, name)

	return bc, nil
}

func (c *Compiler) lvalueTerm(bc *bytecode.ByteCode) error {

	term := c.t.Peek(1)
	if term == "[" {

		c.t.Advance(1)
		ix, err := expressions.Compile(c.t)
		if err != nil {
			return err
		}
		bc.Append(ix)
		if !c.t.IsNext("]") {
			return c.NewError("missing ] on array index")
		}
		bc.Emit0(bytecode.StoreIndex)
		return nil
	}

	if term == "." {
		c.t.Advance(1)
		member := c.t.Next()
		if !tokenizer.IsSymbol(member) {
			return c.NewTokenError("invalid member name")
		}
		bc.Emit(bytecode.Push, member)
		bc.Emit0(bytecode.StoreIndex)
		return nil
	}

	return nil
}
