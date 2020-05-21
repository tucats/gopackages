package compiler

import (
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/tokenizer"
)

const (
	indexLoopType = 1
	rangeLoopType = 2
)

// Loop is a structure that defines a loop type.
type Loop struct {
	Parent *Loop
	Type   int
	// Fixup locations for break or continue statements in a
	// loop. These are the addresses that must be fixed up with
	// the target address.
	breaks    []int
	continues []int
}

// Compiler is a structure defining what we know about the
// compilation
type Compiler struct {
	b         *bytecode.ByteCode
	t         *tokenizer.Tokenizer
	s         *symbols.SymbolTable
	loops     *Loop
	constants []string
}

// Compile starts a compilation unit, and returns a bytecode
// of the compiled material.
func Compile(t *tokenizer.Tokenizer) (*bytecode.ByteCode, error) {

	b := bytecode.New("")
	cInstance := Compiler{b: b, t: t, s: &symbols.SymbolTable{Name: "compile-unit"}, constants: make([]string, 0)}
	c := &cInstance

	c.t.Reset()

	for !c.t.AtEnd() {
		err := c.Statement()
		if err != nil {
			return nil, err
		}
	}

	// Append any symbols created to the bytecode's table
	st := c.Symbols()

	for k, v := range st.Symbols {
		c.b.Symbols.SetAlways(k, v)

	}
	return c.b, nil
}

// StatementEnd returns true when the next token is
// the end-of-statement boundary
func (c *Compiler) StatementEnd() bool {
	next := c.t.Peek(1)

	if next == tokenizer.EndOfTokens {
		return true
	}

	return (next == ";") || (next == "}")
}

// Symbols returns the symbol table map from compilation
func (c *Compiler) Symbols() *symbols.SymbolTable {
	return c.s
}
