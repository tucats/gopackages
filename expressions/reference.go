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
	for e.TokenP < len(e.Tokens)-1 {

		op := e.Tokens[e.TokenP]
		switch op {

		// Structure member reference
		case ".":
			e.TokenP = e.TokenP + 1
			name := e.Tokens[e.TokenP]
			e.TokenP = e.TokenP + 1
			e.b.Emit(bc.Push, name)
			e.b.Emit(bc.Member, nil)

		// Array index reference
		case "[":
			e.TokenP = e.TokenP + 1
			err := e.conditional()
			if err != nil {
				return err
			}
			if e.TokenP > len(e.Tokens)-1 || e.Tokens[e.TokenP] != "]" {
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
