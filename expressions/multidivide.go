package expressions

import (
	bc "github.com/tucats/gopackages/bytecode"
)

// Eval evaluates the parsed expression. This can be called multiple times
// with the same scanned string, but with different symbols.
func (e *Expression) multDivide() error {

	err := e.unary()
	if err != nil {
		return err
	}

	var parsing = true
	for parsing {
		if e.TokenP >= len(e.Tokens) {
			break
		}
		op := e.Tokens[e.TokenP]
		if inList(op, []string{"*", "/", "|"}) {
			e.TokenP = e.TokenP + 1

			err := e.unary()
			if err != nil {
				return err
			}

			switch op {

			case "*":
				e.b.Emit(bc.Mul, 0)

			case "/":
				e.b.Emit(bc.Div, 0)

			case "|":
				e.b.Emit(bc.Or, 0)

			}

		} else {
			parsing = false
		}
	}
	return nil
}
