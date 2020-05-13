package expressions

import (
	bc "github.com/tucats/gopackages/bytecode"
)

// Eval evaluates the parsed expression. This can be called multiple times
// with the same scanned string, but with different symbols.
func (e *Expression) Eval(symbols *bc.SymbolTable) (interface{}, error) {

	// If the compile failed, bail out now.
	if e.err != nil {
		return nil, e.err
	}

	// If the symbol table we're given is unallocated, make one for our use now.
	if symbols == nil {
		symbols = bc.NewSymbolTable("")

	}

	// Add the builtin functions
	AddBuiltins(symbols)

	// Run the generated code to get a result
	ctx := bc.NewContext(symbols, e.b)
	err := ctx.Run()
	if err != nil {
		return nil, err
	}

	return ctx.Pop()
}
