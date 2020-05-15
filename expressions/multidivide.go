package expressions

import (
	bc "github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/tokenizer"
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
		if e.t.AtEnd() {
			break
		}
		op := e.t.Peek(1)
		if tokenizer.InList(op, []string{"*", "/", "|"}) {
			e.t.Advance(1)

			err := e.unary()
			if err != nil {
				return err
			}

			switch op {

			case "*":
				e.b.Emit1(bc.Mul)

			case "/":
				e.b.Emit1(bc.Div)

			case "|":
				e.b.Emit1(bc.Or)

			}

		} else {
			parsing = false
		}
	}
	return nil
}
