package expressions

import (
	"errors"

	bc "github.com/tucats/gopackages/bytecode"
)

// reference parses a structure or array reference
func (e *Expression) reference() error {

	// Parse the atom
	err := e.expressionAtom()
	if err != nil {
		return err
	}

	// is there a trailing structure or array reference?
	for !e.t.AtEnd() {

		op := e.t.Peek(1)
		switch op {

		// Map member reference
		case ".":
			e.t.Advance(1)
			name := e.t.Next()
			e.b.Emit(bc.Push, name)
			e.b.Emit(bc.Member, nil)

		// Array index reference
		case "[":
			e.t.Advance(1)
			err := e.conditional()
			if err != nil {
				return err
			}
			if e.t.Next() != "]" {
				return errors.New("missing ] in array reference")
			}
			e.b.Emit(bc.Index, nil)

		// Nothing else, term is complete
		default:
			return nil
		}
	}
	return nil
}
