package expressions

// Eval evaluates the parsed expression. This can be called multiple times
// with the same scanned string, but with different symbols.
func (e *Expression) Eval(symbols map[string]interface{}) (interface{}, error) {

	// If the compile failed, bail out now.
	if e.err != nil {
		return nil, e.err
	}

	// If the symbol table we're given is unallocated, make one for our use now.
	if symbols == nil {
		symbols = map[string]interface{}{}

	}

	// Add the builtin functions
	AddBuiltins(symbols)

	// Run the generated code to get a result
	err := e.b.Run(symbols)
	if err != nil {
		return nil, err
	}

	return e.b.Pop()
}
