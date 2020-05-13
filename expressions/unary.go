package expressions

import bc "github.com/tucats/gopackages/bytecode"

func (e *Expression) unary() error {

	// Check for unary negation or not before passing into top-level diadic operators.

	t := e.t.Peek()
	switch t {
	case "-":
		e.t.Advance(1)
		err := e.reference()
		if err != nil {
			return err
		}
		e.b.Emit(bc.Negate, 0)
		return nil

	case "!":
		e.t.Advance(1)
		err := e.reference()
		if err != nil {
			return err
		}
		e.b.Emit(bc.Negate, 0)
		return nil

	default:
		return e.reference()

	}
}
