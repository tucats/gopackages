package expressions

import "github.com/tucats/gopackages/bytecode"

// Eval evaluates the parsed expression. This can be called multiple times
// with the same scanned string, but with different symbols.
func (e *Expression) Eval(symbols map[string]interface{}) (interface{}, error) {

	e.TokenP = 0
	e.b = bytecode.New("expression")

	if symbols == nil {
		symbols = map[string]interface{}{}

	}

	AddBuiltins(symbols)

	// Let's check for the special case of an assignment operation
	if len(e.Tokens) > 2 && symbol(e.Tokens[0]) && e.Tokens[1] == ":=" {
		e.TokenP = 2
		v, err := e.conditional(symbols)
		if err != nil {
			return nil, err
		}
		symbols[e.Tokens[0]] = v
		return v, nil
	}
	return e.conditional(symbols)
}
