package compiler

import (
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/tokenizer"
)

// Compiler is a structure defining what we know about the
// compilation
type Compiler struct {
	b *bytecode.ByteCode
	t *tokenizer.Tokenizer
}

// Compile starts a compilation unit, and returns a bytecode
// of the compiled material.
func Compile(t *tokenizer.Tokenizer) (*bytecode.ByteCode, error) {

	b := bytecode.New("compile-unit")
	cInstance := Compiler{b: b, t: t}
	c := &cInstance

	c.t.Reset()

	for !c.t.AtEnd() {
		err := c.Statement()
		if err != nil {
			return nil, err
		}
	}

	return c.b, nil
}

// StatementEnd returns true when the next token is
// the end-of-statement boundary
func (c *Compiler) StatementEnd() bool {
	next := c.t.Peek(1)
	return (next == ";") || (next == "}")
}
