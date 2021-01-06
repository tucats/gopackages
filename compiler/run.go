package compiler

import (
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/tokenizer"
)

// Given a string, compile and execute it immediately.
func RunString(s *symbols.SymbolTable, stmt string) error {
	return Run(s, tokenizer.New(stmt))
}

// Given a token stream, compile and execute it immediately.
func Run(s *symbols.SymbolTable, t *tokenizer.Tokenizer) error {

	c := New()
	bc, err := c.Compile(t)
	if err == nil {
		ctx := bytecode.NewContext(s, bc)
		err = ctx.Run()
	}
	return err
}
