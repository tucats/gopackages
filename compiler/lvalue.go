package compiler

import (
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/tokenizer"
	"github.com/tucats/gopackages/util"
)

// IsLValue peeks ahead to see if this is likely to be an lValue
// object. This is used in cases where the parser might be in an
// otherwise ambiguous state
func (c *Compiler) IsLValue() bool {
	name := c.t.Peek(1)
	if !tokenizer.IsSymbol(name) {
		return false
	}

	// See if it's a reserved word.
	if util.InList(name, tokenizer.ReservedWords...) {
		return false
	}

	next := c.t.Peek(2)
	if next == "." || next == "[" {
		return true
	}

	if util.InList(next, "=", ":=") {
		return true
	}
	return false
}

// LValue compiles the information on the left side of
// an assignment. This information is used later to store the
// data in the named object.
func (c *Compiler) LValue() (*bytecode.ByteCode, error) {

	bc := bytecode.New("lvalue")
	name := c.t.Next()

	if !tokenizer.IsSymbol(name) {
		return nil, c.NewError(InvalidSymbolError, name)
	}
	name = c.Normalize(name)

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

	// Quick optimization; if the name is "_" it just means
	// discard and we can shortcircuit that.

	if name == "_" {
		bc.Emit(bytecode.Drop, 1)
	} else {

		if c.t.Peek(1) == ":=" {
			bc.Emit(bytecode.SymbolCreate, name)
		}

		// Is the last operation in the stack referecing
		// a parent object? If so, convert the last one to
		// a store operation.
		ops := bc.Opcodes()
		opsPos := bc.Mark() - 1
		if opsPos > 0 && ops[opsPos].Opcode == bytecode.LoadIndex {
			ops[opsPos].Opcode = bytecode.StoreIndex
		} else {
			bc.Emit(bytecode.Store, name)
		}
	}
	return bc, nil
}

// lvalueTerm parses secondary lvalue operations (array indexes, or struct member dereferences)
func (c *Compiler) lvalueTerm(bc *bytecode.ByteCode) error {

	term := c.t.Peek(1)
	if term == "[" {

		c.t.Advance(1)
		ix, err := c.Expression()
		if err != nil {
			return err
		}
		bc.Append(ix)
		if !c.t.IsNext("]") {
			return c.NewError(MissingBracketError)
		}
		bc.Emit(bytecode.LoadIndex)
		return nil
	}

	if term == "." {
		c.t.Advance(1)
		member := c.t.Next()
		if !tokenizer.IsSymbol(member) {
			return c.NewError(InvalidSymbolError, member)
		}

		bc.Emit(bytecode.Push, c.Normalize(member))
		bc.Emit(bytecode.LoadIndex)
		return nil
	}

	return nil
}
