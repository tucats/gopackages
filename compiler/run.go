package compiler

import (
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/tokenizer"
)

// Given a string, compile and execute it immediately.
func RunString(name string, s *symbols.SymbolTable, stmt string) error {
	return Run(name, s, tokenizer.New(stmt))
}

// Given a token stream, compile and execute it immediately.
func Run(name string, s *symbols.SymbolTable, t *tokenizer.Tokenizer) error {

	c := New()
	c.ExtensionsEnabled(true)
	bc, err := c.Compile(name, t)
	if err == nil {
		ctx := bytecode.NewContext(s, bc)
		err = ctx.Run()
	}
	return err
}
