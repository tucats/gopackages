package expressions

import "github.com/tucats/gopackages/bytecode"

// Eval evaluates the parsed expression. This can be called multiple times
// with the same scanned string, but with different symbols.
func (e *Expression) Eval(symbols SymbolTable) (interface{}, error) {

	// If the compile failed, bail out now.
	if e.err != nil {
		return nil, e.err
	}

	// If the symbol table we're given is unallocated, make one for our use now.
	if symbols == nil {
		symbols = SymbolTable{}

	}

	// Add the builtin functions
	AddBuiltins(symbols)

	// Run the generated code to get a result
	ctx := bytecode.NewContext(symbols, e.b)
	err := ctx.Run()
	if err != nil {
		return nil, err
	}

	return ctx.Pop()
}
